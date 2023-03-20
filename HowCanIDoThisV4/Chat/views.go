package chat

import (
	"github.com/gorilla/websocket"
)

type MessageForChat struct {
	MessageId    uint64 `json:"messageId"`
	MessageText  string `json:"messageText"`
	UserSenderId uint32 `json:"userSenderId"`
	SendTime     string `json:"sendTime"`
	MessageType  string `json:"messageType"`
	ChatId       uint64 `json:"chatId"`
}
type ViewData struct {
	MesFor      []MessageForChat `json:"mesFor"`
	UserHowSend uint32           `json:"userHowSend"`
}

type TmpStruct struct {
	MessageId    uint64
	Message      string
	MessageType  string
	ChatId       uint64
	SecondUserId uint32
}

type message struct {
	data []byte
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

type ConnectUser struct {
	Websocket *websocket.Conn
	UserID    uint32
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
