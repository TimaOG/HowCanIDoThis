package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

func CheckUserSingIn(w http.ResponseWriter, r *http.Request, email string, password string) bool {

	userData := getUserByEmail(email)

	if userData.Email == "" {
		return false
	}

	if bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password)) != nil {
		return false
	}

	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "userID",
		Value:   fmt.Sprint(userData.Id),
		Expires: expirationTime,
	})
	return true
}

func CheckUserAuth(w http.ResponseWriter, r *http.Request) bool {

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("1")
			return false
		}
		fmt.Println("2")
		return false
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("3")
			return false
		}
		fmt.Println("4")
		return false
	}
	if !tkn.Valid {
		fmt.Println("5")
		return false
	}

	return true
}

func RefreshCookie(w http.ResponseWriter, r *http.Request) {

	c1, err := r.Cookie("userID")
	if err != nil {
		if err == http.ErrNoCookie {
			panic(err)
		}
	}

	userIdStr := c1.Value
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "userID",
		Value:   userIdStr,
		Expires: expirationTime,
	})
}

func UserLogOut(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "userID",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
}
