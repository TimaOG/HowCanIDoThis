package account

import (
	"fmt"
	auth "hcidt/Auth"
	config "hcidt/Config"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func AccountPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	data := getUserByIdFull(uint32(userId))
	t, err := template.ParseFiles("templates/account.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}

func GetThemeNumberPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	themeNumber := getThemeNumber(uint32(userId))
	themeNumberStr := strconv.Itoa(int(themeNumber))
	w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\", \"body\": \"" + themeNumberStr + "\"}")))
}

func ChangeThemeNumberPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	vars := mux.Vars(r)
	themeNumberStr := vars["themeNumber"]
	themeNumber, _ := strconv.Atoi(themeNumberStr)
	changeThemeNumber(uint32(userId), uint8(themeNumber))
	w.Write([]byte("{\"status\": \"ok\"}"))
}

func ChangeActiveStatus(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	changeActiveStatus(uint32(userId))
	w.Write([]byte("{\"status\": \"ok\"}"))
}

func ChangeProfileImg(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	err = r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prevFileName := getPrevFileName(uint32(userId))
	if prevFileName != "" {
		err = os.Remove("./static/upload/avatars/" + prevFileName)
	}
	checkError(err)
	file, h, err := r.FormFile("file")
	checkError(err)
	pathToSave := "./static/upload/avatars/" + userIdStr.Value + h.Filename
	tmpfile, err := os.Create(pathToSave)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, file)
	checkError(err)
	changeProfileImg(uint32(userId), (userIdStr.Value + h.Filename))
	w.Write([]byte("{\"status\": \"ok\", \"body\": \"1\"}"))
}

func NewOrderPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	var returnId uint32
	if r.Method == http.MethodPost {
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			panic(err)
		}
		orderName := r.FormValue("orderName")
		orderDiscribtion := r.Form.Get("orderDiscribction")
		orderDeadline := r.FormValue("orderDeadline")
		_, err = time.Parse("2006-01-02", orderDeadline)
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderPrice, err := strconv.Atoi(r.FormValue("orderPrice"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderUrgency, err := strconv.Atoi(r.FormValue("orderUrgency"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderWorkType, err := strconv.Atoi(r.FormValue("orderWorkType"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderCategory, err := strconv.Atoi(r.FormValue("orderCategoryFirst"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderCategorySecond, err := strconv.Atoi(r.FormValue("orderCategorySecond"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		orderTags := r.FormValue("orderTags")
		var fn string
		if r.MultipartForm != nil && r.MultipartForm.File != nil {
			file, h, err := r.FormFile("fileTZ")
			if err != nil {
				fn = ""
			} else {
				pathToSave := "./static/upload/tz/" + userIDstr.Value + h.Filename
				fn = userIDstr.Value + h.Filename
				_, err := os.Stat(pathToSave)
				if os.IsNotExist(err) {
					tmpfile, err := os.Create(pathToSave)
					checkError(err)
					_, err = io.Copy(tmpfile, file)
					defer tmpfile.Close()
					checkError(err)
				} else {
					new_names := strings.Split(h.Filename, ".")
					fn = new_names[0] + "N." + new_names[1]
					npathToSave := "./static/upload/offerSkin/" + userIDstr.Value + fn
					tmpfile, err := os.Create(npathToSave)
					checkError(err)
					_, err = io.Copy(tmpfile, file)
					defer tmpfile.Close()
					checkError(err)
				}
			}
		} else {
			fn = ""
		}
		var order = Order{orderName, uint32(userID), orderDiscribtion, uint32(orderPrice), orderDeadline,
			uint8(orderUrgency), uint8(orderWorkType), uint8(orderCategory), uint8(orderCategorySecond), orderTags, true, fn}

		returnId = addOrder(order)

	}
	rId := strconv.Itoa(int(returnId))

	w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\", \"body\": \"Предложение создано\", \"lastId\": \"" + rId + "\"}")))

}

func NewOfferPage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	var returnId uint32
	var fn string
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		offerName := r.FormValue("offerName")
		offerDiscribtion := r.FormValue("offerDiscribtion")
		daysToComplite, err := strconv.Atoi(r.FormValue("offerDaysToComplite"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		offerPrice, _ := strconv.Atoi(r.FormValue("offerPrice"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"2\"}"))
			return
		}
		offerWorkType, err := strconv.Atoi(r.FormValue("offerWorkType"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"3\"}"))
			return
		}
		offerCategory, err := strconv.Atoi(r.FormValue("offerCategoryFirst"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"4\"}"))
			return
		}
		offerCategorySecond, err := strconv.Atoi(r.FormValue("offerCategorySecond"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"5\"}"))
			return
		}
		offerTags := r.FormValue("offerTags")
		file, h, err := r.FormFile("cover")
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"6\"}"))
			return
		} else {
			pathToSave := "./static/upload/offerSkin/" + userIDstr.Value + h.Filename
			fn = userIDstr.Value + h.Filename
			_, err := os.Stat(pathToSave)
			if os.IsNotExist(err) {
				tmpfile, err := os.Create(pathToSave)
				checkError(err)
				_, err = io.Copy(tmpfile, file)
				defer tmpfile.Close()
				checkError(err)
			} else {
				new_names := strings.Split(h.Filename, ".")
				fn = new_names[0] + "N." + new_names[1]
				npathToSave := "./static/upload/offerSkin/" + userIDstr.Value + fn
				tmpfile, err := os.Create(npathToSave)
				checkError(err)
				_, err = io.Copy(tmpfile, file)
				defer tmpfile.Close()
				fn = userIDstr.Value + fn
				checkError(err)
			}
		}
		var offer = Offer{offerName, uint32(userID), offerDiscribtion, uint32(offerPrice), uint8(daysToComplite),
			uint8(offerWorkType), uint8(offerCategory), uint8(offerCategorySecond), offerTags, true, fn}
		returnId = addOffer(offer)

	}
	rId := strconv.Itoa(int(returnId))
	w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\", \"body\": \"Предложение создано\", \"lastId\": \"" + rId + "\", \"fn\": \"" + fn + "\"}")))
}

func ChangeActiveStatusForOffer(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	if r.Method == http.MethodPost {
		r.ParseForm()
		offerIdStr := r.Form.Get("value")
		offerId, _ := strconv.Atoi(offerIdStr)
		changeActiveStatusForOffer(uint32(userId), uint32(offerId))
		w.Write([]byte("{\"status\": \"ok\", \"body\": \"Статус предложения изменен\"}"))
		return
	}
	w.Write([]byte("{\"status\": \"fail\", \"body\": \"ERROR\"}"))
}

func ChangeAccountInfo(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	if r.Method == http.MethodPost {
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		userPassword := getUserPasswordById(uint32(userID))
		email := r.FormValue("newEmail")
		oldPassword := r.FormValue("oldPassword")
		newPassword := r.FormValue("newPassword")
		newPassword2 := r.FormValue("newPassword2")
		if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(oldPassword)) == nil {
			if checkForEmailInSystem(email) && email != "" {
				saveAccountEmail(uint32(userID), email)
			}
			if (newPassword == newPassword2) && (newPassword != "" && newPassword2 != "") {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				saveNewPassword(uint32(userID), hashedPassword)
			}
		} else {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"Неверный пароль\"}"))
			return
		}
	}
	w.Write([]byte("{\"status\": \"ok\", \"body\": \"Изменения сохранены\"}"))
}

func ChangeOfferGetOldData(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	vars := mux.Vars(r)
	offerIdStr := vars["offer"]
	offerId, _ := strconv.Atoi(offerIdStr)
	result := getOfferById(uint32(offerId))
	if string(result) != "no" {
		w.Write(result)
		return
	}
	w.Write([]byte("{\"status\": \"fail\", \"body\": \"ERROR\"}"))
}
func ChangeOffer(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	var fn string
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		offerIdStr := r.FormValue("offerId")
		offerId, _ := strconv.Atoi(offerIdStr)
		offerName := r.FormValue("offerName")
		offerDiscribtion := r.FormValue("offerDiscribtion")
		daysToComplite, err := strconv.Atoi(r.FormValue("offerDaysToComplite"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"1\"}"))
			return
		}
		offerPrice, _ := strconv.Atoi(r.FormValue("offerPrice"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"2\"}"))
			return
		}
		offerWorkType, err := strconv.Atoi(r.FormValue("offerWorkType"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"3\"}"))
			return
		}
		offerCategory, err := strconv.Atoi(r.FormValue("offerCategoryFirst"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"4\"}"))
			return
		}
		offerCategorySecond, err := strconv.Atoi(r.FormValue("offerCategorySecond"))
		if err != nil {
			w.Write([]byte("{\"status\": \"fail\", \"body\": \"5\"}"))
			return
		}
		offerTags := r.FormValue("offerTags")
		file, h, err := r.FormFile("cover")
		if err != nil {
			fn = ""
		} else {
			pathToSave := "./static/upload/offerSkin/" + userIDstr.Value + h.Filename
			fn = userIDstr.Value + h.Filename
			_, err := os.Stat(pathToSave)
			if os.IsNotExist(err) {
				tmpfile, err := os.Create(pathToSave)
				checkError(err)
				_, err = io.Copy(tmpfile, file)
				defer tmpfile.Close()
				checkError(err)
			} else {
				new_names := strings.Split(h.Filename, ".")
				fn = new_names[0] + "N." + new_names[1]
				npathToSave := "./static/upload/offerSkin/" + userIDstr.Value + fn
				tmpfile, err := os.Create(npathToSave)
				checkError(err)
				_, err = io.Copy(tmpfile, file)
				defer tmpfile.Close()
				checkError(err)
			}
		}
		var offer = OfferById{"ok", offerName, offerDiscribtion, uint32(offerPrice), uint8(daysToComplite),
			uint8(offerWorkType), uint8(offerCategory), uint8(offerCategorySecond), offerTags}
		changeOffer(offer, uint32(offerId), fn, uint32(userID))

	}
	w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\", \"body\": \"Предложение измененно\", \"fn\": \"" + fn + "\"}")))
}

func DeleteOffer(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	if r.Method == http.MethodPost {
		r.ParseForm()
		offerIdStr := r.Form.Get("value")
		offerId, _ := strconv.Atoi(offerIdStr)
		fn := deleteOffer(uint32(userId), uint32(offerId))
		pathToDel := "./static/upload/offerSkin/" + fn
		_, err := os.Stat(pathToDel)
		if !os.IsNotExist(err) {
			err = os.Remove(pathToDel)
			checkError(err)
		}
		w.Write([]byte("{\"status\": \"ok\", \"body\": \"Предложение удалено\"}"))
		return
	}
	w.Write([]byte("{\"status\": \"fail\", \"body\": \"ERROR\"}"))
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		//http.Redirect(w, r, fmt.Sprintf(config.SiteAdress+"login?target=account"), http.StatusSeeOther)
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	auth.RefreshCookie(w, r)
	userIdStr, err := r.Cookie("userID")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(userIdStr.Value)
	if r.Method == http.MethodPost {
		r.ParseForm()
		orderIdStr := r.Form.Get("value")
		orderId, _ := strconv.Atoi(orderIdStr)
		fn := deleteOrder(uint32(userId), uint32(orderId))
		if fn != "" {
			pathToDel := "./static/upload/tz/" + fn
			_, err := os.Stat(pathToDel)
			if !os.IsNotExist(err) {
				err = os.Remove(pathToDel)
				checkError(err)
			}
			err = os.Remove(pathToDel)
			checkError(err)
		}
		checkError(err)
		w.Write([]byte("{\"status\": \"ok\", \"body\": \"Предложение удалено\"}"))
		return
	}
	w.Write([]byte("{\"status\": \"fail\", \"body\": \"ERROR\"}"))
}

func AccountLogout(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
		return
	}
	auth.UserLogOut(w, r)

	http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/login?target=account"), http.StatusSeeOther)
}

func AccountDeletePage(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckUserAuth(w, r) {
		w.Write([]byte("{\"status\": \"notAuth\", \"body\": \"1\"}"))
		return
	}
	if r.Method == http.MethodPost {
		userIDstr, _ := r.Cookie("userID")
		userID, _ := strconv.Atoi(userIDstr.Value)
		userPassword := getUserPasswordById(uint32(userID))
		password := r.FormValue("password")
		if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)) != nil {
			http.Redirect(w, r, fmt.Sprintf(config.Http+r.Host+"/account"), http.StatusSeeOther)
			return
		}
		auth.UserLogOut(w, r)
		deleteAccountById(uint32(userID))
	}
	http.Redirect(w, r, fmt.Sprintf(r.Host+"/account"), http.StatusSeeOther)
}

func AccountWalletPage(w http.ResponseWriter, r *http.Request) {

}
func AccountWalletFillPage(w http.ResponseWriter, r *http.Request) {

}
func AccountBringOutPage(w http.ResponseWriter, r *http.Request) {

}
