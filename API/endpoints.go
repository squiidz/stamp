package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	//"github.com/squiidz/stamp/module/logger"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"html/template"
)

type Users struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Create   time.Time `json:"create"`
	Update   time.Time `json:"update"`
	Friends  []string  `json:"friends"`
	Pending  []string  `json:"pending"` // Friend Request Pending
}

type Friend struct {
	Username    string `json:"username"`
	LastMessage string `json:"lastMessage"`
}

type Location struct {
	Longitude float64
	Latitude  float64
}

type Message struct {
	From     Users
	To       []string
	Message  string
	Create   time.Time
	Position Location
	Picture  []byte
}

const (
	MongoServerAddr         = "localhost"
	RedisServerAddr         = "192.168.0.104"
	TemplateFolder          = "/template"
	StaticFolder            = "/static"
	SessionTTL      float64 = 5.00
)

var (
	MongoSession, err = mgo.Dial(MongoServerAddr)

	MDB  = MongoSession.DB("message")
	MCol = MDB.C("new")
	MSav = MDB.C("save")

	UDB  = MongoSession.DB("account")
	UCol = UDB.C("user")

	TemplatesLocation = map[string]string{}
	Templates         = template.New("main")
	Store             = sessions.NewCookieStore([]byte("BigBangBazooka"))
)

// Check User Position , and return Message if they exist for the current location
func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	sessionC, _ := Store.Get(req, "sessionCookie")
	//logger.CheckErr(err, "COOKIE DOESN'T EXIST")
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
	//logger.CheckErr(err, "CANNOT ADD FRIEND")

	username := CookieValue(req, "sessionCookie")

	err = UpdateFriendList(username, string(friendData))
	if err != nil {
		//logger.CheckErr(err, "USER "+string(friendData)+" DOESN'T EXISTS !")
		fmt.Fprintf(rw, "non-valid")
	}
}
