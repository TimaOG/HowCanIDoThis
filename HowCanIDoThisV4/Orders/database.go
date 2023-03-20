package orders

import (
	"database/sql"
	"fmt"
	config "hcidt/Config"
	"strconv"
	"strings"
	"unicode/utf8"
)

var dbURL = config.DbURL

func ordersListWithoutSorting() ViewData {
	db, _ := sql.Open("postgres", dbURL)
	sqlRequest := `SELECT t1.id, t1.ordername, t1.fkUserOwner, t2.username, t1.discribtion ,t1.price, t1.deadline, t1.urgency, t1.workType, t1.tags FROM Orders as t1 JOIN Users as t2 on t1.FkUserOwner = t2.id WHERE t1.isActive = true LIMIT 20;`
	res, err := db.Query(sqlRequest)
	checkError(err)
	orders := []Order{}
	for res.Next() {
		order := Order{}
		var tmpDisc string
		err = res.Scan(&order.Id, &order.Name, &order.FkUserOwner, &order.FkUserOwnerName, &tmpDisc, &order.Price, &order.Deadline, &order.Urgency, &order.WorkType, &order.Tags)
		checkError(err)
		tmpDeadline := strings.Split(order.Deadline, "T")
		order.Deadline = tmpDeadline[0]
		if utf8.RuneCountInString(tmpDisc) > 20 {
			order.Discribtion = tmpDisc[0:20]
		} else {
			order.Discribtion = tmpDisc
		}
		orders = append(orders, order)
	}
	sqlRequest = `select * from TaskTypeFirst`
	res2, err := db.Query(sqlRequest)
	checkError(err)
	taskTypesFirst := []TaskTypeFirst{}
	for res2.Next() {
		taskTypeFirst := TaskTypeFirst{}
		err = res2.Scan(&taskTypeFirst.Id, &taskTypeFirst.Name)
		checkError(err)
		taskTypesFirst = append(taskTypesFirst, taskTypeFirst)
	}
	sqlRequest = `select * from TaskTypeSecond`
	res3, err := db.Query(sqlRequest)
	checkError(err)
	taskTypesSecond := []TaskTypeSecond{}
	for res3.Next() {
		taskTypeSecond := TaskTypeSecond{}
		err = res3.Scan(&taskTypeSecond.Id, &taskTypeSecond.Name, &taskTypeSecond.FkFirstType)
		checkError(err)
		taskTypesSecond = append(taskTypesSecond, taskTypeSecond)
	}
	return ViewData{orders, taskTypesFirst, taskTypesSecond}
}

func ordersListWithSorting(sortFields SortFields) []Order {
	db, _ := sql.Open("postgres", dbURL)
	sqlRequest := getSelectWithFilters(sortFields)
	res, err := db.Query(sqlRequest)
	checkError(err)
	orders := []Order{}
	for res.Next() {
		order := Order{}
		var tmpDisc string
		err = res.Scan(&order.Id, &order.Name, &order.FkUserOwner, &order.FkUserOwnerName, &tmpDisc, &order.Price, &order.Deadline, &order.Urgency, &order.WorkType, &order.Tags)
		checkError(err)
		tmpDeadline := strings.Split(order.Deadline, "T")
		order.Deadline = tmpDeadline[0]
		if utf8.RuneCountInString(tmpDisc) > 20 {
			order.Discribtion = tmpDisc[0:20]
		} else {
			order.Discribtion = tmpDisc
		}
		orders = append(orders, order)
	}
	return orders
}

func getSelectWithFilters(sortFields SortFields) string {
	sqlRequest := `SELECT t1.id, t1.ordername, t1.fkUserOwner, t2.username, t1.discribtion ,t1.price, t1.deadline, t1.urgency, t1.workType, t1.tags FROM Orders as t1 JOIN Users as t2 on t1.FkUserOwner = t2.id WHERE t1.isActive = true `
	if sortFields.TaskType != "0" {
		sqlRequest += "and t1.taskTypeFirst = " + fmt.Sprint(sortFields.TaskType) + " "
	}
	if sortFields.SecondType != "0" {
		sqlRequest += "and t1.taskTypeSecond = " + fmt.Sprint(sortFields.SecondType) + " "
	}
	if sortFields.PriceDown != "0" || sortFields.PriceUp != "0" {
		if sortFields.PriceDown != "" {
			sqlRequest += "AND t1.price > " + fmt.Sprint(sortFields.PriceDown) + " "
		} else if sortFields.PriceUp != "" {
			sqlRequest += "AND t1.price < " + fmt.Sprint(sortFields.PriceUp) + " "
		}
	}
	if sortFields.WorkType != "0" {
		sqlRequest += "AND t1.workType = " + fmt.Sprint(sortFields.WorkType) + " "
	}
	if sortFields.Urgency != "0" {
		sqlRequest += "AND t1.urgency = " + fmt.Sprint(sortFields.WorkType) + " "
	}
	offset, _ := strconv.Atoi(sortFields.Offset)
	sqlRequest += "ORDER BY t1.price " + sortFields.OrderBy + " LIMIT 20 OFFSET " + fmt.Sprint(offset*20)
	return sqlRequest
}

func getOrderById(id uint32) OrderById {
	sqlRequest := `SELECT t1.id, t1.ordername, t1.fkUserOwner, t2.username, t1.discribtion ,t1.price, t1.deadline, t1.urgency, t1.workType, t1.tags, t1.tzPath FROM Orders as t1 JOIN Users as t2 on t1.FkUserOwner = t2.id WHERE t1.id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var order OrderById
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&order.Id, &order.Name, &order.FkUserOwner, &order.FkUserOwnerName, &order.Discribtion, &order.Price, &order.Deadline, &order.Urgency, &order.WorkType, &order.Tags, &order.TzPath)
	if err == sql.ErrNoRows {
		return OrderById{}
	}
	return order
}

func makeChatAndMath(userId uint32, userOwnerID uint32, taskId uint32) {
	sqlRequest := "INSERT INTO Matches (fkUserWhoDo, fkUserTaskOwner, fkWhatTaskId) VALUES(?, ?, ?) RETURNING id;"
	db, _ := sql.Open("postgres", dbURL)
	var tmpId uint64
	err := db.QueryRow(sqlRequest, userId, userOwnerID, taskId).Scan(&tmpId)
	checkError(err)
	sqlRequest = "INSERT INTO Chats (fkUserSecond, fkUserFirst, fkWhatTaskId) VALUES(?, ?, ?) RETURNING id;"
	defer db.Close() //Доделать
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
