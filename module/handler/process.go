package handler

import (
	"encoding/json"
	"github.com/squiidz/stamp/module/logger"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

const (
	SmallZone float64 = 0.0005200
	MedZone   float64 = 0.0012000
	BigZone   float64 = 1.0005200
)

func CheckUser(user Users) bool {
	err = UCol.Find(bson.M{"username": user.Username, "password": user.Password}).One(&user)
	if err != nil {
		return false
	} else {
		return true
	}
}

func FindUser(username string) *Users {
	user := Users{}
	err := UCol.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		logger.CheckErr(err, "CANNOT FIND CONNECTED USER")
	}
	return &user
}

func CheckMessage(username *string, loc *Location, rw *http.ResponseWriter) {
	var messCheck = make(chan bool, 10)
	m := []Message{}

	err := MCol.Find(bson.M{"to": username}).All(&m)
	logger.CheckErr(err, "CANNOT FIND MESSAGE")

	// Encode to Json messages found
	enco := json.NewEncoder(*rw)
	//log.Println("MESSAGE COUNT : ", len(m))

	// Loop over all find messages
	for _, mess := range m {
		go PositionValid(&mess, loc, messCheck)
		if <-messCheck {
			enco.Encode(&mess)
			MCol.Update(bson.M{"message": mess.Message}, bson.M{"$pull": bson.M{"to": *username}})
			go CheckIfEmpty(&mess)
		}
	}
}

func PositionValid(message *Message, location *Location, check chan bool) {
	var zone float64 = BigZone

	if (location.Latitude-message.Position.Latitude) < zone || (location.Latitude-message.Position.Latitude) > (zone-zone*2) && (location.Latitude-message.Position.Latitude) < zone {
		if (location.Longitude-message.Position.Longitude) < zone || (location.Longitude-message.Position.Longitude) > (zone-zone*2) && (location.Longitude-message.Position.Longitude) < zone {
			log.Println("TRUE")
			check <- true
		} else {
			check <- false
		}
	} else {
		check <- false
	}
}

func CheckIfEmpty(m *Message) {
	err := MCol.Remove(bson.M{"from.username": m.From.Username, "message": m.Message, "to": bson.M{"$size": 0}})
	if err != nil {
		logger.CheckErr(err, "CANT DELETE MESSAGE")
	}
}
