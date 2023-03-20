package chat

import (
	"encoding/json"
	"fmt"
	auth "hcidt/Auth"
	config "hcidt/Config"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func (H *hub) Run() {
	for {
		select {
		case s := <-H.register:
			connections := H.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				H.rooms[s.room] = connections
			}
			H.rooms[s.room][s.conn] = true
		case s := <-H.unregister:
			connections := H.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(H.rooms, s.room)
					}
				}
			}
		case m := <-H.broadcast:
			connections := H.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(H.rooms, m.room)
					}
				}
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		H.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{msg, s.room}
		H.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump(roomId uint64, userId uint64) {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			var tmp TmpStruct
			err := json.Unmarshal(message, &tmp)
			idToSend, _ := strconv.Atoi(tmp.Id)
			dt := time.Now()
			mesToSendJson := MessageForChat{tmp.Message, uint32(idToSend), dt.Format("01-02-2006 15:04:05"), tmp.FileName}
			json_data, err := json.Marshal(mesToSendJson)
			if err != nil {
				log.Fatal(err)
			}
			if err := c.write(websocket.TextMessage, json_data); err != nil {
				return
			} else {
				go saveMessage(roomId, tmp.Message, uint64(idToSend), tmp.FileName)
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, roomId string, userId uint32) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	H.register <- s
	roomIdInt, _ := strconv.Atoi(roomId)
	go s.writePump(uint64(roomIdInt), uint64(userId))
	go s.readPump()
}

func ChatPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
		return
	}
	auth.RefreshCookie(w, r)
	userIDstr, _ := r.Cookie("userID")
	userId, _ := strconv.Atoi(userIDstr.Value)
	data := getChadIds(uint32(userId))
	cData := ChatsData{data, uint32(userId)}
	t, err := template.ParseFiles("templates/chats.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, cData)
}

func ChatSaveFile(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	auth.RefreshCookie(w, r)
	var fn string
	userIDstr, _ := r.Cookie("userID")
	file, h, err := r.FormFile("document")
	if err != nil {
		w.Write([]byte("{\"status\": \"fail\", \"body\": \"6\"}"))
		return
	} else {
		pathToSave := "./static/upload/messeges/" + userIDstr.Value + h.Filename
		fn = userIDstr.Value + h.Filename
		_, err := os.Stat(pathToSave)
		if os.IsNotExist(err) {
			tmpfile, err := os.Create(pathToSave)
			checkError(err)
			_, err = io.Copy(tmpfile, file)
			defer tmpfile.Close()
			checkError(err)
		} else {
			new_names := strings.Split(h.Filename, ".")
			fn = new_names[0] + "||N." + new_names[1]
			npathToSave := "./static/upload/offerSkin/" + userIDstr.Value + fn
			tmpfile, err := os.Create(npathToSave)
			checkError(err)
			_, err = io.Copy(tmpfile, file)
			defer tmpfile.Close()
			fn = userIDstr.Value + fn
			checkError(err)
		}
	}
}

func ChatPageId(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=chat"), http.StatusSeeOther)
		return
	}
	auth.RefreshCookie(w, r)
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	roomId, _ := strconv.Atoi(vars["id"])
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=chat/"+vars["id"]), http.StatusSeeOther)
		return
	}
	messageData := getMessage(uint64(roomId))
	userIDstr, _ := r.Cookie("userID")
	userId, _ := strconv.Atoi(userIDstr.Value)
	data := ViewData{messageData, uint32(userId)}
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	checkError(err)
	w.Write(jsonData)
}

func ChatPageWsId(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=chat"), http.StatusSeeOther)
		return
	}
	auth.RefreshCookie(w, r)
	vars := mux.Vars(r)
	roomId := vars["id"]
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=chat/"+vars["id"]), http.StatusSeeOther)
		return
	}
	userIDstr, _ := r.Cookie("userID")
	userId, _ := strconv.Atoi(userIDstr.Value)

	serveWs(w, r, roomId, uint32(userId))
}
