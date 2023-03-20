package auth

import (
	"database/sql"
	config "hcidt/Config"

	_ "github.com/lib/pq"
)

var dbURL = config.DbURL

func getUserByEmail(email string) User {
	sqlRequest := `SELECT id, userpassword, email FROM Users WHERE email = $1;`
	db, _ := sql.Open("postgres", dbURL)
	var user User
	res := db.QueryRow(sqlRequest, email)
	err := res.Scan(&user.Id, &user.Password, &user.Email)
	if err == sql.ErrNoRows {
		return User{}
	}
	checkError(err)
	defer db.Close()
	return user
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
