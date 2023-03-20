package other

import (
	"fmt"
	auth "hcidt/Auth"
	config "hcidt/Config"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func StartPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/startPage.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	data := ""
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if auth.CheckUserSingIn(w, r, email, password) {
			target := r.URL.Query().Get("target")
			if len(target) == 0 {
				target = "account"
			}
			http.Redirect(w, r, (config.Http + r.Host + "/" + target), http.StatusSeeOther)
			return
		} else {
			data = "Неверный логин или пароль"
		}
	}
	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	data := ""
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		name := r.FormValue("name")

		createOrNot := true
		if password != password2 {
			data = "пароли не равны"
			createOrNot = false
		} else if checkForUserInSystem(email) >= 1 {
			data = "Пользователь уже есть в системе"
			createOrNot = false
		}
		if createOrNot {
			addUser(name, email, password)
			http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login"), http.StatusSeeOther)
			return
		}
	}
	t, err := template.ParseFiles("templates/registration.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)

}
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	auth.UserLogOut(w, r)
	http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login"), http.StatusSeeOther)

}
func HelpPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/helpPage.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func UserInfoPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["userName"])
	data := getUserByIdShort(uint32(userId))
	t, err := template.ParseFiles("templates/user.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
