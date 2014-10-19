package handler

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// Check User Position , and return Message if they exist for the current location
func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	m := []Message{}
	sessionC, err := Store.Get(req, "sessionCookie")

	log.Println("[+] Session Name Value : ", sessionC.Values["name"].(string)) // Print Cookie Name Value for debuggin

	if err != nil {
		log.Println("COOKIE DOESN'T EXIST")
	}
	username := sessionC.Values["name"].(string)
	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	err = MCol.Find(bson.M{"to.username": username}).All(&m)

	if err != nil {
		log.Println("CANNOT FIND MESSAGE")
		log.Println(err) // Return nothing if no messages
	}

	// Encode to Json messages found
	enco := json.NewEncoder(rw)
	for _, mess := range m {
		if PositionValid(&mess, &loc) {
			log.Println(username, " #", len(m))
			enco.Encode(&mess)
			MCol.Remove(bson.M{"to.username": username, "message": mess.Message})
			/*
				log.Println("{Message Position}")
				log.Println("[Lat] : ",mess.Latitude, "[Long] : ", mess.Longitude)
				log.Println("{User Postion}")
				log.Println("[Lat] : ", loc.Latitude, "[Long] : ", loc.Longitude)
			*/
		}
	}
}

// Insert New Message to Database
func PlaceHandler(rw http.ResponseWriter, req *http.Request) {
	message := Message{}
	data := json.NewDecoder(req.Body)
	data.Decode(&message)

	log.Println("New Message for : ", message.To[0].Username)

	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}
	MCol.Insert(&message)
}
