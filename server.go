package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/go-zoo/claw"
	"github.com/go-zoo/claw/mw"

	hand "github.com/squiidz/stamp/module/handler"
	mid "github.com/squiidz/stamp/module/middle"
)

var port string
var mux = bone.New()

func init() {
	portFlag := flag.String("port", "8080", "-port [8080]")
	flag.Parse()
	port = ":" + *portFlag
}

func main() {
	cl := claw.New(mw.Logger)

	mux.Handle("/", cl.Use(hand.LoginHandler))
	mux.Handle("/register", cl.Use(hand.RegisterHandler))
	mux.Handle("/home/:id", mid.AuthMiddle(cl.Use(hand.IndexHandler), "sessionCookie"))
	mux.Handle("/profil", mid.AuthMiddle(cl.Use(hand.ProfilHandler), "sessionCookie"))
	mux.Post("/addfriend", cl.Use(hand.AddFriendHandler))
	mux.Post("/save", cl.Use(hand.SaveHandler))
	mux.Post("/insert", cl.Use(hand.InsertMessageHandler))
	mux.Post("/location", mid.AuthMiddle(cl.Use(hand.LocationHandler), "sessionCookie"))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("SERVER RUNNIG ON PORT %s \n", port)

	http.ListenAndServe(port, mux)
}
