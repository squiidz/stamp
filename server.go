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
	portFlag := flag.String("port", ":80", "-port [8080]")
	flag.Parse()
	port = ":" + *portFlag
}

func main() {
	mux.HandleFunc("/", hand.LoginHandler)
	mux.HandleFunc("/home", hand.IndexHandler)
	mux.HandleFunc("/watch", hand.WatchHandler)
	mux.HandleFunc("/place", hand.PlaceHandler)
	mux.HandleFunc("/location", hand.LocationHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("SERVER RUNNIG ON PORT %s \n", port)

	http.ListenAndServe(port, mux)
}
