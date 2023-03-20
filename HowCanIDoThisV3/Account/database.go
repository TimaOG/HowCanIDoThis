package account

import (
	"database/sql"
	"encoding/json"
	config "hcidt/Config"
	"strings"
	"unicode/utf8"

	_ "github.com/lib/pq"
)

var dbURL = config.DbURL

func getUserByIdFull(id uint32) ViewData {
	sqlRequest := `SELECT name, email, discribtion, profileImg, rating, balance ,isPremiumUser, isActiveUser, 
	historyCount, responsibility, doneOnTime, answerSpead, registrationDate
	FROM users WHERE id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var user UserInfo
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&user.Name, &user.Email, &user.Discribtion, &user.ProfileImg, &user.Rating, &user.Balance, &user.IsPremiumUser, &user.IsActiveUser,
		&user.HistotyCount, &user.Responsibility, &user.DoneOnTime, &user.AnswerSpead, &user.RegistrationDate)
	tmpDate := strings.Split(user.RegistrationDate, "T")
	user.RegistrationDate = tmpDate[0]
	if err == sql.ErrNoRows {
		return ViewData{}
	}
	checkError(err)
	sqlRequest = `SELECT id, name, discribtion, price FROM Orders WHERE fkUserOwner = $1;`
	res1, err := db.Query(sqlRequest, id)
	checkError(err)
	orders := []UserOrders{}
	for res1.Next() {
		order := UserOrders{}
		var tmpDisc string
		err := res1.Scan(&order.Id, &order.Name, &tmpDisc, &order.Price)
		if utf8.RuneCountInString(tmpDisc) > 20 {
			order.Discribtion = tmpDisc[0:20]
		} else {
			order.Discribtion = tmpDisc
		}
		checkError(err)
		orders = append(orders, order)
	}
	sqlRequest = `SELECT id, name, price, isActive, coverPath FROM Offers WHERE fkUserOwner = $1;`
	res2, err := db.Query(sqlRequest, id)
	checkError(err)
	offers := []UserOffers{}
	for res2.Next() {
		offer := UserOffers{}
		err = res2.Scan(&offer.Id, &offer.Name, &offer.Price, &offer.IsActive, &offer.CoverPath)
		checkError(err)
		offers = append(offers, offer)
	}
	sqlRequest = `SELECT * FROM getchatsbyuserid($1)`
	res3, err := db.Query(sqlRequest, id)
	checkError(err)
	chats := []UserChats{}
	for res3.Next() {
		chat := UserChats{}
		err = res3.Scan(&chat.UserChatId, &chat.UserSecondId, &chat.UserSecondName, &chat.UserSecondImg, &chat.LastMessageText, &chat.LastMessageTime)
		checkError(err)
		chats = append(chats, chat)
	}
	sqlRequest = `select * from TaskTypeFirst`
	res4, err := db.Query(sqlRequest)
	checkError(err)
	taskTypesFirst := []TaskTypeFirst{}
	for res4.Next() {
		taskTypeFirst := TaskTypeFirst{}
		err = res4.Scan(&taskTypeFirst.Id, &taskTypeFirst.Name)
		checkError(err)
		taskTypesFirst = append(taskTypesFirst, taskTypeFirst)
	}
	sqlRequest = `select * from TaskTypeSecond`
	res5, err := db.Query(sqlRequest)
	checkError(err)
	taskTypesSecond := []TaskTypeSecond{}
	for res5.Next() {
		taskTypeSecond := TaskTypeSecond{}
		err = res5.Scan(&taskTypeSecond.Id, &taskTypeSecond.Name, &taskTypeSecond.FkFirstType)
		checkError(err)
		taskTypesSecond = append(taskTypesSecond, taskTypeSecond)
	}
	sqlRequest = `SELECT t4.id, t4.name, t4.discribtion, t4.price, t2.name, t2.id, t3.name, t3.id, 
	t1.startTime, t4.deadline FROM Mathes as t1 LEFT JOIN Users as t2 on t1.fkUserTaskOwner = t2.id 
	LEFT JOIN Users as t3 on t1.fkUserWhoDo = t3.id LEFT JOIN Orders as t4 on t1.fkWhatTaskId = t4.id
	WHERE t1.fkUserTaskOwner = $1 AND t1.userOwnerConfim = true AND  t1.userWhoDoConfim = true AND isOrder = true;`
	res6, err := db.Query(sqlRequest, id)
	checkError(err)
	ordersDN := []UserOrdersDoingNow{}
	for res6.Next() {
		orderDN := UserOrdersDoingNow{}
		err = res6.Scan(&orderDN.Id, &orderDN.Name, &orderDN.Discribtion, &orderDN.Price, &orderDN.ExecutorName,
			&orderDN.ExecutorId, &orderDN.CustomerName, orderDN.CustomerId, &orderDN.StartTime, &orderDN.ExpectedEndTime)
		checkError(err)
		ordersDN = append(ordersDN, orderDN)
	}
	sqlRequest = `SELECT t4.id, t4.name, t4.price, t2.name, t2.id, t3.name, t3.id, 
	t1.startTime, t4.deadline FROM Mathes as t1 LEFT JOIN Users as t2 on t1.fkUserTaskOwner = t2.id 
	LEFT JOIN Users as t3 on t1.fkUserWhoDo = t3.id LEFT JOIN Orders as t4 on t1.fkWhatTaskId = t4.id
	WHERE t1.fkUserTaskOwner = $1 AND t1.userOwnerConfim = true AND  t1.userWhoDoConfim = true AND isOrder = false;`
	res7, err := db.Query(sqlRequest, id)
	checkError(err)
	offersDN := []UserOffersDoingNow{}
	for res7.Next() {
		offerDN := UserOffersDoingNow{}
		err = res7.Scan(&offerDN.Id, &offerDN.Name, &offerDN.Price, &offerDN.ExecutorName, &offerDN.ExecutorId,
			&offerDN.CustomerName, &offerDN.CustomerId, &offerDN.StartTime, &offerDN.ExpectedEndTime)
		checkError(err)
		offersDN = append(offersDN, offerDN)
	}
	defer db.Close()
	return ViewData{user, orders, offers, ordersDN, offersDN, chats, taskTypesFirst, taskTypesSecond}
}

func changeActiveStatus(id uint32) {
	sqlRequest := `CALL changeactivity($1)`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, id)
	checkError(err)
	defer db.Close()
}
func changeActiveStatusForOffer(id uint32, offerId uint32) {
	sqlRequest := `CALL changeactivityforoffer($1, $2)`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, id, offerId)
	checkError(err)
	defer db.Close()
}
func changeProfileImg(id uint32, pathToSave string) {
	sqlRequest := `UPDATE Users SET profileImg = $1 WHERE id = $2`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, pathToSave, id)
	checkError(err)
	defer db.Close()
}

func getPrevFileName(id uint32) string {
	sqlRequest := `SELECT ProfileImg FROM Users WHERE id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var prevFN string
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&prevFN)
	checkError(err)
	return prevFN
}

func getOfferById(id uint32) []byte {
	sqlRequest := `SELECT name, discribtion, price, daysToComplite, workType, taskTypeFirst, taskTypeSecond, tags FROM Offers WHERE id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var offer OfferById
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&offer.Name, &offer.Discribtion, &offer.Price, &offer.DaysToComplite, &offer.WorkType, &offer.OrderCategory, &offer.OrderCategorySecond, &offer.Tags)
	if err == sql.ErrNoRows {
		return []byte("no")
	}
	offer.Status = "ok"
	returnRes, _ := json.Marshal(offer)
	return returnRes
}

func addOrder(order Order) uint32 {
	sqlRequest := `INSERT INTO Orders (name, fkUserOwner, discribtion, price, deadline, urgency, workType,
		taskTypeFirst, taskTypeSecond, tags, isActive, tzPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	db, _ := sql.Open("postgres", dbURL)
	var returnId uint32
	res := db.QueryRow(sqlRequest, order.Name, order.FkUserOwner, order.Discribtion, order.Price, order.Deadline,
		order.Urgency, order.WorkType, order.OrderCategory, order.OrderCategorySecond, order.Tags, order.IsActive, order.TzPath)
	err := res.Scan(&returnId)
	checkError(err)
	defer db.Close()
	return returnId
}
func addOffer(offer Offer) uint32 {
	sqlRequest := `INSERT INTO Offers (name, fkUserOwner, discribtion, price, daysToComplite, workType,
		taskTypeFirst, taskTypeSecond, tags, isActive, coverPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	db, _ := sql.Open("postgres", dbURL)
	var returnId uint32
	res := db.QueryRow(sqlRequest, offer.Name, offer.FkUserOwner, offer.Discribtion, offer.Price, offer.DaysToComplite,
		offer.WorkType, offer.OrderCategory, offer.OrderCategorySecond, offer.Tags, offer.IsActive, offer.CoverPath)
	err := res.Scan(&returnId)
	checkError(err)
	defer db.Close()
	return returnId
}
func changeOffer(offer OfferById, offerId uint32, fn string, userId uint32) {
	sqlRequest := ""
	db, _ := sql.Open("postgres", dbURL)
	if fn == "" {
		sqlRequest = `UPDATE Offers SET name = $1, discribtion = $2, price = $3, daysToComplite = $4, workType = $5,
		taskTypeFirst = $6, taskTypeSecond = $7, tags = $8 WHERE id = $9 AND FkUserOwner = $10`
		_, err := db.Exec(sqlRequest, offer.Name, offer.Discribtion, offer.Price, offer.DaysToComplite,
			offer.WorkType, offer.OrderCategory, offer.OrderCategorySecond, offer.Tags, offerId, userId)
		checkError(err)
	} else {
		sqlRequest = `UPDATE Offers SET name = $1, discribtion = $2, price = $3, daysToComplite = $4, workType = $5,
		taskTypeFirst = $6, taskTypeSecond = $7, tags = $8, coverPath = $9 WHERE id = $10 AND FkUserOwner = $11`
		_, err := db.Exec(sqlRequest, offer.Name, offer.Discribtion, offer.Price, offer.DaysToComplite,
			offer.WorkType, offer.OrderCategory, offer.OrderCategorySecond, offer.Tags, fn, offerId, userId)
		checkError(err)
	}
	defer db.Close()
}
func deleteOffer(id uint32, offerId uint32) string {
	sqlRequest := `DELETE FROM Offers WHERE id = $1 AND fkUserOwner = $2 RETURNING coverPath`
	db, _ := sql.Open("postgres", dbURL)
	var fn string
	res := db.QueryRow(sqlRequest, offerId, id)
	res.Scan(&fn)
	defer db.Close()
	return fn
}
func deleteOrder(id uint32, orderId uint32) string {
	sqlRequest := `DELETE FROM Orders WHERE id = $1 AND fkUserOwner = $2 RETURNING tzPath`
	db, _ := sql.Open("postgres", dbURL)
	var fn string
	res := db.QueryRow(sqlRequest, orderId, id)
	res.Scan(&fn)
	defer db.Close()
	return fn
}

func getUserPasswordById(id uint32) string {
	sqlRequest := `SELECT password FROM Users WHERE id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var userPassword string
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&userPassword)
	checkError(err)
	return userPassword
}

func checkForEmailInSystem(email string) bool {
	sqlRequest := "SELECT COUNT(id) FROM Users WHERE email = $1;"
	db, _ := sql.Open("postgres", dbURL)
	var answ int
	res := db.QueryRow(sqlRequest, email)
	err := res.Scan(&answ)
	checkError(err)
	defer db.Close()
	if answ == 0 {
		return true
	}
	return false
}

func deleteAccountById(id uint32) {
	sqlRequest := `DELETE FROM Users WHERE id = $1`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, id)
	checkError(err)
	defer db.Close()
}

func saveAccountEmail(id uint32, email string) {
	sqlRequest := `UPDATE Users SET email = $1 WHERE id = $2`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, email, id)
	checkError(err)
	defer db.Close()
}

func saveNewPassword(id uint32, newPassword []byte) {
	sqlRequest := `UPDATE Users SET password = $1 WHERE id = $2`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, newPassword, id)
	checkError(err)
	defer db.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
