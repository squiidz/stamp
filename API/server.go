package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/go-zoo/claw"
	"github.com/go-zoo/claw/mw"
)

var (
	port        string
	fakeFriends = [2]Friend{}
)

func init() {
	portFlag := flag.String("port", "8080", "-port [8080]")
	flag.Parse()
	port = ":" + *portFlag

	fakeFriends[0] = Friend{
		Username:    "Alex",
		LastMessage: "Pra Pra Pra !",
	}
}

func main() {
	muxx := bone.New()
	cw := claw.New(mw.Logger)

	muxx.Get("/friends", cw.Use(FriendsHandler))
	muxx.Get("/profil", cw.Use(ProfilHandler))
	muxx.Get("/message", cw.Use(MessageHandler))
	//muxx.Handle("/", cw.Use(LoginHandler))
	//muxx.Handle("/register", cw.Use(RegisterHandler))
	//muxx.Handle("/home/:id", mid.AuthMiddle(cw.Use(IndexHandler), "sessionCookie"))
	//muxx.Handle("/profil", mid.AuthMiddle(cw.Use(ProfilHandler), "sessionCookie"))

	muxx.Post("/addfriend", cw.Use(AddFriendHandler))
	muxx.Post("/save", cw.Use(SaveHandler))
	muxx.Post("/insert", cw.Use(InsertMessageHandler))
	//muxx.Handle("/location", cw.Use(LocationHandler))

  http.ListenAndServe(port, muxx)
}

func FriendsHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Request friends from %s\n", req.RemoteAddr)
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	enco := json.NewEncoder(rw)
	enco.Encode(fakeFriends)
}

func ProfilHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Request Profil Info from %s\n", req.RemoteAddr)
	rw.Header().Add("Access-Control-Allow-Origin", "*")

	user, err := GetUser("squiidz")
	if err != nil {
		fmt.Println(err)
	}

	enco := json.NewEncoder(rw)
	enco.Encode(user)
}