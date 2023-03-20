package offers

import (
	"encoding/json"
	"fmt"
	auth "hcidt/Auth"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func OffersListPage(w http.ResponseWriter, r *http.Request) {
	data := offersListWithoutSorting()
	t, err := template.ParseFiles("templates/marketOffer.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
func OffersListPageSort(w http.ResponseWriter, r *http.Request) {
	var jsonData []byte
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		var sortFields SortFields
		//var filters map[string]string
		err := json.NewDecoder(r.Body).Decode(&sortFields)
		fmt.Println(sortFields)
		checkError(err)
		data := offersListWithSorting(sortFields)
		jsonData, err = json.Marshal(data)
		checkError(err)
	}
	w.Write(jsonData)
}
func OfferPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIdStr := vars["offer"]
	taskId, _ := strconv.Atoi(taskIdStr)
	data := getOfferById(uint32(taskId))
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	checkError(err)
	w.Write(jsonData)
}

func OfferMessagePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	vars := mux.Vars(r)
	taskIdStr := vars["offer"]
	if r.Method == http.MethodPost {
		r.ParseForm()
		taskId, _ := strconv.Atoi(taskIdStr)
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		messageText := r.Form.Get("value")
		makeChat(uint32(userID), uint32(taskId), messageText)
		w.Write([]byte("{\"status\": \"ok\", \"body\": \"Сообщение отправлено\"}"))
		return
	}
	w.Write([]byte("{\"status\": \"fail\", \"body\": \"ERROR\"}"))
}
