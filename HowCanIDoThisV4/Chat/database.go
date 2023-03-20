package chat

import (
	"database/sql"
	config "hcidt/Config"
	"log"
)

var dbURL = config.DbURL

func getChadIds(userId uint32) []Chat {
	sqlRequest := `SELECT * FROM getchatsbyuserid($1)`
	db, _ := sql.Open("postgres", dbURL)
	res, err := db.Query(sqlRequest, userId)
	if err != nil {
		panic(err.Error())
	}
	chats := []Chat{}
	for res.Next() {
		chat := Chat{}
		err := res.Scan(&chat.Id, &chat.SecondUserId, &chat.SecondUserName, &chat.SecondUserImg, &chat.LastMessage, &chat.LastSendTime)
		if err != nil {
			log.Fatal(err)
		}
		chats = append(chats, chat)
	}
	defer db.Close()
	return chats

}

func saveMessage(chatId uint64, message string, userId uint64, messageType string) {
	if messageType == "text" || messageType == "file" || messageType == "mathes" {
		var sqlRequest string
		db, _ := sql.Open("postgres", dbURL)
		sqlRequest = "CALL saveMessage($1, $2, $3, $4)"
		_, err := db.Query(sqlRequest, chatId, message, userId, messageType)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
	}
}

func getMessage(chatId uint64) []MessageForChat {
	sqlRequest := "SELECT messageText, fkUserId, sendTime, messagetype FROM Messages WHERE fkChatId = $1"
	db, _ := sql.Open("postgres", dbURL)
	res, err := db.Query(sqlRequest, chatId)
	if err != nil {
		panic(err.Error())
	}
	chatMessages := make([]MessageForChat, 0)
	for res.Next() {
		var mfc MessageForChat
		err := res.Scan(&mfc.MessageText, &mfc.UserSenderId, &mfc.SendTime, &mfc.MessageType)
		if err != nil {
			log.Fatal(err)
		}
		chatMessages = append(chatMessages, mfc)
	}
	defer res.Close()
	defer db.Close()
	return chatMessages
}

func makeMatch(userId uint32, userWhoDo uint32, taskId uint32, isOrder bool) {
	sqlRequest := "CALL makeMatch($1, $2, $3, $4);"
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, userId, userWhoDo, taskId, isOrder)
	checkError(err)
}

func confimMatch(userId uint32, taskId uint32) {
	sqlRequest := "CALL confimMatch($1, $2);"
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, userId, taskId)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
