package main

import (
	"flag"
	"fmt"
	hand "github.com/squiidz/stamp/module/handler"
	mid "github.com/squiidz/stamp/module/middle"
	"net/http"
)

var port string
var mux = http.NewServeMux()

func init() {
	portFlag := flag.String("port", "8080", "-port [8080]")
	flag.Parse()
	port = ":" + *portFlag
}

func main() {
	mux.Handle("/", http.HandlerFunc(hand.LoginHandler))
	mux.Handle("/register", http.HandlerFunc(hand.RegisterHandler))
	mux.Handle("/home", mid.AuthMiddle(http.HandlerFunc(hand.IndexHandler), "sessionCookie"))
	mux.Handle("/profil", mid.AuthMiddle(http.HandlerFunc(hand.ProfilHandler), "sessionCookie"))
	mux.Handle("/addfriend", http.HandlerFunc(hand.AddFriendHandler))
	mux.Handle("/save", http.HandlerFunc(hand.SaveHandler))
	mux.Handle("/insert", http.HandlerFunc(hand.InsertMessageHandler))
	mux.Handle("/location", mid.AuthMiddle(http.HandlerFunc(hand.LocationHandler), "sessionCookie"))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("SERVER RUNNIG ON PORT %s \n", port)

	http.ListenAndServe(port, mux)
}
