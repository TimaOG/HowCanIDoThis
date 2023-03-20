package other

import (
	"database/sql"
	config "hcidt/Config"

	"golang.org/x/crypto/bcrypt"
)

var dbURL = config.DbURL

func addUser(userName string, userEmail string, userPassword string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	sqlRequest := `INSERT INTO Users (username, email, userpassword) VALUES ($1, $2, $3)`
	db, _ := sql.Open("postgres", dbURL)
	_, err := db.Exec(sqlRequest, userName, userEmail, hashedPassword)
	checkError(err)
	defer db.Close()
}

func checkForUserInSystem(email string) int {
	sqlRequest := "SELECT COUNT(id) FROM Users WHERE email = $1;"
	db, _ := sql.Open("postgres", dbURL)
	var answ int
	res := db.QueryRow(sqlRequest, email)
	err := res.Scan(&answ)
	checkError(err)
	defer db.Close()
	return answ
}

func getUserByIdShort(id uint32) User {
	sqlRequest := `SELECT id, username, discribtion, profileImg, rating, isPremiumUser, isActiveUser 
	historyCount, responsibility, doneOnTime, answerSpead, registrationDate
	FROM Users WHERE id = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var user User
	res := db.QueryRow(sqlRequest, id)
	err := res.Scan(&user.Id, &user.Name, &user.Discribtion, &user.ProfileImg, &user.Rating, &user.IsPremiumUser, &user.IsActiveUser, &user.HistotyCount, &user.Responsibility, &user.DoneOnTime, &user.AnswerSpead, &user.RegistrationDate)
	if err == sql.ErrNoRows {
		return User{}
	}
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return user
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
