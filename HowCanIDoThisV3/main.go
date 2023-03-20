package main

import (
	"fmt"
	account "hcidt/Account"
	chat "hcidt/Chat"
	offers "hcidt/Offers"
	orders "hcidt/Orders"
	other "hcidt/Other"
	"net/http"

	"github.com/gorilla/mux"
)

func startServerAndAddHandles() {
	r := mux.NewRouter()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", other.StartPage)
	r.HandleFunc("/login", other.LoginPage)
	r.HandleFunc("/registration", other.RegistrationPage)
	r.HandleFunc("/logout", other.LogoutPage)

	r.HandleFunc("/help", other.HelpPage)

	r.HandleFunc("/user/{userName}", other.UserInfoPage)

	r.HandleFunc("/account", account.AccountPage)
	r.HandleFunc("/account/changeActiveStatus", account.ChangeActiveStatus)
	r.HandleFunc("/account/changeProfileImg", account.ChangeProfileImg)
	r.HandleFunc("/account/newOrder", account.NewOrderPage)
	r.HandleFunc("/account/newOffer", account.NewOfferPage)
	r.HandleFunc("/account/changeActiveStatusForOffer", account.ChangeActiveStatusForOffer)
	r.HandleFunc("/account/changeSettings", account.ChangeAccountInfo)
	r.HandleFunc("/account/getOfferInfo/{offer}", account.ChangeOfferGetOldData)
	r.HandleFunc("/account/changeOffer", account.ChangeOffer)
	r.HandleFunc("/account/deleteOffer", account.DeleteOffer)
	r.HandleFunc("/account/deleteOrder", account.DeleteOrder)
	r.HandleFunc("/account/logout", account.AccountLogout)
	r.HandleFunc("/account/delete", account.AccountDeletePage)
	r.HandleFunc("/account/wallet", account.AccountWalletPage)
	r.HandleFunc("/account/wallet/fill", account.AccountWalletFillPage)
	r.HandleFunc("/account/wallet/BringOut", account.AccountBringOutPage)

	r.HandleFunc("/chat", chat.ChatPage)
	r.HandleFunc("/chat/{id:[0-9]+}", chat.ChatPageId)
	r.HandleFunc("/chat/ws/{id:[0-9]+}", chat.ChatPageWsId)
	r.HandleFunc("/chat/saveChatFile", chat.ChatSaveFile)

	r.HandleFunc("/orders", orders.OrdersListPage)
	r.HandleFunc("/orders/sort", orders.OrdersListPageSort)
	r.HandleFunc("/orders/{order}", orders.OrderPage)
	r.HandleFunc("/orders/{order}/message", orders.OrderMessagePage)

	r.HandleFunc("/offers", offers.OffersListPage)
	r.HandleFunc("/offers/sort", offers.OffersListPageSort)
	r.HandleFunc("/offers/{offer}", offers.OfferPage)
	r.HandleFunc("/offers/{offer}/message", offers.OfferMessagePage)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func main() {
	go chat.H.Run()
	fmt.Println("http://localhost:8080")
	startServerAndAddHandles()
}
