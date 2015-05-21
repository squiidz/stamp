package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/squiidz/stamp/Web/module/logger"
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
	fmt.Println(message)
	message.Create = time.Now()
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

// Add friend
func AddFriendHandler(rw http.ResponseWriter, req *http.Request) {
	friendData, err := ioutil.ReadAll(req.Body)
	logger.CheckErr(err, "CANNOT ADD FRIEND")

	username := CookieValue(req, "sessionCookie")

	err = UpdateFriendList(username, string(friendData))
	if err != nil {
		logger.CheckErr(err, "USER "+string(friendData)+" DOESN'T EXISTS !")
		fmt.Fprintf(rw, "non-valid")
	}
}
