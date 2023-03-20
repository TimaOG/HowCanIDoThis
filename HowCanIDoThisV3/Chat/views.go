package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

type MessageForChat struct {
	MessageText  string `json:"messageText"`
	UserSenderId uint32 `json:"userSenderId"`
	SendTime     string `json:"sendTime"`
	FilePath     string `json:"filePath"`
}
type ViewData struct {
	MesFor      []MessageForChat `json:"mesFor"`
	UserHowSend uint32           `json:"userHowSend"`
}

type TmpStruct struct {
	Id       string
	Message  string
	FileName string
}

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn *connection
	room string
}

type MessegeStructur struct {
	Texts      []string
	UserSender string
}

type ChatsData struct {
	Chats  []Chat
	UserId uint32
}

type Chat struct {
	Id             uint64
	SecondUserId   uint32
	SecondUserName string
	SecondUserImg  string
	LastMessage    string
	LastSendTime   string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}
