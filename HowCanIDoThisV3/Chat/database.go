package chat

import (
	"database/sql"
	config "hcidt/Config"
	"log"
	"time"
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

func saveMessage(chatId uint64, message string, userId uint64, fileName string) {
	dt := time.Now()
	var sqlRequest string
	var err error
	db, _ := sql.Open("postgres", dbURL)
	if fileName == "" {
		sqlRequest = "INSERT INTO Messages (fkChatId, messageText, fkUserId, sendTime) VALUES ($1, $2, $3, $4)"
		_, err = db.Query(sqlRequest, chatId, message, userId, dt.Format("01-02-2006 15:04:05"))
	} else {
		sqlRequest = "INSERT INTO Messages (fkChatId, messageText, fkUserId, sendTime, filepath) VALUES ($1, $2, $3, $4, $5)"
		_, err = db.Query(sqlRequest, chatId, message, userId, dt.Format("01-02-2006 15:04:05"), fileName)
	}
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

}

func getMessage(chatId uint64) []MessageForChat {
	sqlRequest := "SELECT messageText, fkUserId, sendTime, filepath FROM Messages WHERE fkChatId = $1"
	db, _ := sql.Open("postgres", dbURL)
	res, err := db.Query(sqlRequest, chatId)
	if err != nil {
		panic(err.Error())
	}
	chatMessages := make([]MessageForChat, 0)
	for res.Next() {
		var tmpFilePath sql.NullString
		var mfc MessageForChat
		err := res.Scan(&mfc.MessageText, &mfc.UserSenderId, &mfc.SendTime, &tmpFilePath)
		if tmpFilePath.Valid {
			mfc.FilePath = tmpFilePath.String
		}
		if err != nil {
			log.Fatal(err)
		}
		chatMessages = append(chatMessages, mfc)
	}
	defer res.Close()
	defer db.Close()
	return chatMessages
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
