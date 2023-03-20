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

func newConnectUser(ws *websocket.Conn, userId uint32) *ConnectUser {
	return &ConnectUser{
		Websocket: ws,
		UserID:    userId,
	}
}

// var users = make(map[ConnectUser]uint32)
var users = make(map[uint32]ConnectUser)

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, userId uint32) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()
	var socketClient *ConnectUser = newConnectUser(ws, userId)
	//users[*socketClient] = userId
	users[userId] = *socketClient
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Ws disconnect waiting", err.Error())
			delete(users, userId)
			return
		}
		var tmp TmpStruct
		err = json.Unmarshal(message, &tmp)
		dt := time.Now()
		mesToSendJson := MessageForChat{tmp.MessageId, tmp.Message, uint32(userId), dt.Format("01-02-2006 15:04:05"), tmp.MessageType, tmp.ChatId}
		json_data, err := json.Marshal(mesToSendJson)
		if err != nil {
			log.Fatal(err)
		}

		if err = users[userId].Websocket.WriteMessage(messageType, json_data); err != nil {
			log.Println("Cloud not send Message to ", users[userId].UserID, err.Error())
		} else {
			go saveMessage(tmp.ChatId, tmp.Message, uint64(userId), tmp.MessageType)
		}
		secondUserId := tmp.SecondUserId
		userToSend, ok := users[uint32(secondUserId)]
		if ok {
			if err = userToSend.Websocket.WriteMessage(messageType, json_data); err != nil {
				log.Println("Cloud not send Message to ", users[userId].UserID, err.Error())
			}
		}
	}
}

func ChatPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
		return
	}
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
	userIDstr, _ := r.Cookie("userID")
	userId, _ := strconv.Atoi(userIDstr.Value)

	serveWs(w, r, uint32(userId))
}

func MakeMatchPage(w http.ResponseWriter, r *http.Request) {
	// if !auth.CheckUserAuth(w, r) {
	// 	http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=chat"), http.StatusSeeOther)
	// 	return
	// }
	// auth.RefreshCookie(w, r)
	// userIDstr, _ := r.Cookie("userID")
	// userId, _ := strconv.Atoi(userIDstr.Value)
	// taskIdStr := r.FormValue("taskId")
	// taskId, _ := strconv.Atoi(taskIdStr)
	// userWhoDoStr := r.FormValue("userWhoDo")
	// userWhoDo, _ := strconv.Atoi(userWhoDoStr)
	// makeMatch(uint32(userId), uint32(userWhoDo), uint32(taskId))
}
