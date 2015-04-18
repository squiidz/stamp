package handler

import (
	"github.com/gorilla/sessions"
	//"github.com/squiidz/stamp/module/logger"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"html/template"
	//"net/http"
	"time"
)

type Users struct {
	Username string
	Email    string
	Password string
	Create   time.Time
	Update   time.Time
	Friends  []string
	Pending  []string // Friend Request Pending
}

type Friend struct {
	Name string
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
	Picture  string
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

func init() {
	TemplatesLocation["profil"] = "template/profil.html"
	TemplatesLocation["login"] = "template/login.html"
	TemplatesLocation["register"] = "template/register.html"
	TemplatesLocation["index"] = "template/index.html"

	Templates.Delims("((", "))")
	for _, value := range TemplatesLocation {
		Templates.ParseFiles(value)
	}
}
