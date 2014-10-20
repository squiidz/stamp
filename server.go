package main

import (
	"flag"
	"fmt"
	hand "github.com/squiidz/stamp/module/handler"
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
	mux.HandleFunc("/", hand.LoginHandler)
	mux.HandleFunc("/register", hand.RegisterHandler)
	mux.HandleFunc("/home", hand.IndexHandler)
	mux.HandleFunc("/profil", hand.ProfilHandler)
	mux.HandleFunc("/save", hand.SaveHandler)
	mux.HandleFunc("/insert", hand.InsertMessageHandler)
	mux.HandleFunc("/location", hand.LocationHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("SERVER RUNNIG ON PORT %s \n", port)

	http.ListenAndServe(port, mux)
}
