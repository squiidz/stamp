package handler

import (
	"encoding/json"
	"github.com/squiidz/stamp/module/logger"
	"net/http"
)

// Check User Position , and return Message if they exist for the current location
func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	sessionC, err := Store.Get(req, "sessionCookie")
	logger.CheckErr(err, "COOKIE DOESN'T EXIST")
	username := sessionC.Values["name"].(string)

	//log.Println("[+] Session Name Value : ", username) // Print Cookie Name Value for debuggin
	CheckMessage(&username, &loc, &rw)
}

// Insert New Message to Database
func InsertMessageHandler(rw http.ResponseWriter, req *http.Request) {
	message := Message{}
	data := json.NewDecoder(req.Body)
	data.Decode(&message)

	go MCol.Insert(&message)
}

// Save Messages
func SaveHandler(rw http.ResponseWriter, req *http.Request) {
	message := Message{}
	decode := json.NewDecoder(req.Body)
	decode.Decode(&message)
	// Insert Message in Persistant Database
	go MSav.Insert(&message)
}
