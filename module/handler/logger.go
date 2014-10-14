package handler

import (
	"log"
	"net/http"
)

func SimpleLog(req *http.Request) {
	log.Println(req.Host, req.RemoteAddr, req.URL.RequestURI())
}
