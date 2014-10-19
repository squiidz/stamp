package logger

import (
	"io/ioutil"
	"log"
	"net/http"
)

func SimpleLog(req *http.Request) {
	log.Println(req.Host, req.RemoteAddr, req.URL.RequestURI())
}

func CheckReqBody(req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
}
