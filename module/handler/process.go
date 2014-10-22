package handler

import (
	"encoding/json"
	"github.com/squiidz/stamp/module/logger"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var (
	Zone = map[string]float64{}
)

func init() {
	Zone["small"] = float64(0.0005200)
	Zone["medium"] = float64(0.0012000)
	Zone["big"] = float64(1.0005200)
}

func CheckUser(user Users) bool {
	err = UCol.Find(bson.M{"username": user.Username, "password": user.Password}).One(&user)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (u *Users) FindUser() {
	//user := Users{}
	err := UCol.Find(bson.M{"username": u.Username}).One(u)
	log.Println(u.Username)
	if err != nil {
		logger.CheckErr(err, "CANNOT FIND CONNECTED USER")
	}
	//return &user
}

func CheckMessage(username *string, loc *Location, rw *http.ResponseWriter) {
	var messCheck = make(chan bool, 10)
	m := []Message{}

	err := MCol.Find(bson.M{"to": username}).All(&m)
	logger.CheckErr(err, "CANNOT FIND MESSAGE")

	// Encode to Json messages found
	enco := json.NewEncoder(*rw)
	log.Println("MESSAGE COUNT : ", len(m))

	// Loop over all find messages
	for _, mess := range m {
		go PositionValid(&mess, loc, messCheck)
		if <-messCheck {
			enco.Encode(&mess)
			mess.UpdateMessage(*username)
		}
	}
}

func (m *Message) UpdateMessage(username string) {
	err := MCol.Update(bson.M{"message": m.Message}, bson.M{"$pull": bson.M{"to": username}})
	logger.CheckErr(err, "Cannot Update Message")
	MCol.Remove(bson.M{"message": m.Message, "to": bson.M{"$size": "0"}})
}

func UpdateFriendList(username string, friend string) {
	if friend != "" {
		err := UCol.Update(bson.M{"username": username}, bson.M{"$push": bson.M{"friends": friend}})
		logger.CheckErr(err, "CANNOT UPDATE USER FRIEND LIST")
	} else {
		log.Println("Friend Name = ", friend)
	}
}

func PositionValid(message *Message, location *Location, check chan bool) {
	zone := Zone["big"]

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

func CookieValue(req *http.Request, cookie string) string {
	session, _ := Store.Get(req, cookie)
	data := session.Values["name"].(string)
	return data
}
