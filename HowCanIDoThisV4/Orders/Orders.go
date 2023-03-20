package orders

import (
	"encoding/json"
	"fmt"
	auth "hcidt/Auth"
	config "hcidt/Config"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func OrdersListPage(w http.ResponseWriter, r *http.Request) {
	data := ordersListWithoutSorting()
	t, err := template.ParseFiles("templates/marketOrder.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
func OrdersListPageSort(w http.ResponseWriter, r *http.Request) {
	var jsonData []byte
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		var sortFields SortFields
		err := json.NewDecoder(r.Body).Decode(&sortFields)
		checkError(err)
		data := ordersListWithSorting(sortFields)
		jsonData, err = json.Marshal(data)
		checkError(err)
	}
	w.Write(jsonData)
}
func OrderPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	taskIdStr := vars["order"]
	taskId, _ := strconv.Atoi(taskIdStr)
	data := getOrderById(uint32(taskId))
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	checkError(err)
	w.Write(jsonData)
}

func OrderMessagePage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"login"), http.StatusSeeOther)
		return
	}
	auth.RefreshCookie(w, r)
	vars := mux.Vars(r)
	taskIdStr := vars["offer"]
	if r.Method == http.MethodPost {
		userOwnerID, _ := strconv.Atoi(r.FormValue("priceDown"))
		taskId, _ := strconv.Atoi(taskIdStr)
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		makeChatAndMath(uint32(userID), uint32(userOwnerID), uint32(taskId))
	}
	http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"offers/"+taskIdStr), http.StatusSeeOther)
}
