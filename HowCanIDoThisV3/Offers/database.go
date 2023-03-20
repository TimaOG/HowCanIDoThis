package offers

import (
	"database/sql"
	"fmt"
	config "hcidt/Config"
	"strconv"
)

var dbURL = config.DbURL

func offersListWithoutSorting() ViewData {
	db, _ := sql.Open("postgres", dbURL)
	sqlRequest := `SELECT t1.id, t1.name, t1.fkUserOwner, t2.name, t2.isPremiumUser ,t1.coverPath, t1.price, t1.daysToComplite, t1.workType, t1.tags, t1.rating, t1.historyCount FROM Offers as t1 LEFT JOIN Users as t2 on t1.fkUserOwner = t2.id WHERE t1.isActive = true LIMIT 20;`
	res, err := db.Query(sqlRequest)
	checkError(err)
	offers := []Offer{}
	var tmpRating sql.NullInt16
	var tmpHistoryCount sql.NullInt32
	for res.Next() {
		offer := Offer{}
		err = res.Scan(&offer.Id, &offer.Name, &offer.FkUserOwner, &offer.UserOwnerName, &offer.IsPremiumUserOwner, &offer.CoverPath, &offer.Price, &offer.DaysToComplite, &offer.WorkType, &offer.Tags, &tmpRating, &tmpHistoryCount)
		checkError(err)
		if tmpRating.Valid {
			offer.Rating = uint8(tmpRating.Int16)
		} else {
			offer.Rating = 0
		}
		if tmpHistoryCount.Valid {
			offer.HistoryCount = uint32(tmpHistoryCount.Int32)
		} else {
			offer.HistoryCount = 0
		}
		offers = append(offers, offer)
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
	return ViewData{offers, taskTypesFirst, taskTypesSecond}
}

func offersListWithSorting(sortFields SortFields) []Offer {
	db, _ := sql.Open("postgres", dbURL)
	sqlRequest := getSelectWithFilters(sortFields)
	res, err := db.Query(sqlRequest)
	checkError(err)
	offers := []Offer{}
	var tmpRating sql.NullInt16
	var tmpHistoryCount sql.NullInt32
	for res.Next() {
		offer := Offer{}
		err = res.Scan(&offer.Id, &offer.Name, &offer.FkUserOwner, &offer.UserOwnerName, &offer.IsPremiumUserOwner, &offer.CoverPath, &offer.Price, &offer.DaysToComplite, &offer.WorkType, &offer.Tags, &tmpRating, &tmpHistoryCount)
		checkError(err)
		if tmpRating.Valid {
			offer.Rating = uint8(tmpRating.Int16)
		} else {
			offer.Rating = 0
		}
		if tmpHistoryCount.Valid {
			offer.HistoryCount = uint32(tmpHistoryCount.Int32)
		} else {
			offer.HistoryCount = 0
		}
		offers = append(offers, offer)
	}
	return offers
}

func getSelectWithFilters(sortFields SortFields) string {
	sqlRequest := `SELECT t1.id, t1.name, t1.fkUserOwner, t2.name, t2.isPremiumUser ,t1.coverPath, t1.price, t1.daysToComplite, t1.workType, t1.tags, t1.rating, t1.historyCount FROM Offers as t1 JOIN Users as t2 on t1.fkUserOwner = t2.id WHERE t1.isActive = true `
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
	} else {
		sqlRequest += "AND t1.workType is not null "
	}
	if sortFields.Rating != "0" {
		sqlRequest += "AND t1.rating > " + fmt.Sprint(sortFields.Rating) + " "
	}
	offset, _ := strconv.Atoi(sortFields.Offset)
	sqlRequest += "ORDER BY t1.price " + sortFields.OrderBy + " LIMIT 20 OFFSET " + fmt.Sprint(offset*20)
	fmt.Println(sqlRequest)
	return sqlRequest
}

func getOfferById(id uint32) OfferById {
	sqlRequest := `SELECT t1.id, t1.name, t1.discribtion ,t1.fkUserOwner, t2.name, t2.isPremiumUser ,t1.coverPath, t1.price, t1.daysToComplite, t1.workType, t1.tags, t1.rating, t1.historyCount FROM Offers as t1 LEFT JOIN Users as t2 on t1.fkUserOwner = t2.id WHERE t1.id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var offer OfferById
	res := db.QueryRow(sqlRequest, id)
	var tmpRating sql.NullInt16
	var tmpHistoryCount sql.NullInt32
	err := res.Scan(&offer.Id, &offer.Name, &offer.Discribtion, &offer.FkUserOwner, &offer.UserOwnerName, &offer.IsPremiumUserOwner, &offer.CoverPath, &offer.Price, &offer.DaysToComplite, &offer.WorkType, &offer.Tags, &tmpRating, &tmpHistoryCount)
	if tmpRating.Valid {
		offer.Rating = uint8(tmpRating.Int16)
	} else {
		offer.Rating = 0
	}
	if tmpHistoryCount.Valid {
		offer.HistoryCount = uint32(tmpHistoryCount.Int32)
	} else {
		offer.HistoryCount = 0
	}
	if err == sql.ErrNoRows {
		return OfferById{}
	}
	return offer
}

func makeChat(userId uint32, offerId uint32, messageText string) {
	sqlRequest := "CALL makechat($1, $2, $3)"
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, userId, offerId, messageText)
	checkError(err)
	defer db.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
