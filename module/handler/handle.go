package handler

import (
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/squiidz/stamp/module/logger"
)

type Users struct {
	Username string
	Password string
	Create   time.Time
	Update   time.Time
	Friends  []Friend
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
	To       []Users
	Message  string
	Position Location
}

const (
	MongoServerAddr         = "192.168.0.104"
	RedisServerAddr         = "192.168.0.104"
	TemplateFolder          = "/template"
	StaticFolder            = "/static"
	SessionTTL      float64 = 5.00
)

var (
	MongoSession, err = mgo.Dial(MongoServerAddr)

	MDB  = MongoSession.DB("message")
	MCol = MDB.C("new")

	UDB  = MongoSession.DB("account")
	UCol = UDB.C("user")

	TemplatesLocation = map[string]string{}
	Templates         = template.New("main")
	Store             = sessions.NewCookieStore([]byte("BigBangBazooka"))
)

func init() {
	TemplatesLocation["login"] = "template/login.html"
	TemplatesLocation["index"] = "template/index.html"

	Templates.Delims("((", "))")
	for _, value := range TemplatesLocation {
		Templates.ParseFiles(value)
	}
}

func LoginHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		Templates.ExecuteTemplate(rw, "login.html", nil)
	} else if req.Method == "POST" {
		user := Users{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}
		if CheckUser(user) {
			session, _ := Store.Get(req, "sessionCookie")
			session.Values["name"] = user.Username
			session.Save(req, rw)
			http.Redirect(rw, req, "/home", http.StatusFound)
		} else {
			http.Redirect(rw, req, "/", http.StatusFound)
		}
	}
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	session, _ := Store.Get(req, "sessionCookie")
	data := session.Values["name"].(string)
	Templates.ExecuteTemplate(rw, "index.html", data)
	logger.SimpleLog(req)
}
