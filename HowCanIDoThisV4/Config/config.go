package config

import "fmt"

const (
	User     = "postgres"
	Password = "123"
	DBip     = "127.0.0.1"
	DBPort   = 5432
	DBName   = "test"
)

var DbURL = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBip, DBPort, User, Password, DBName)

const Http = "http://"
