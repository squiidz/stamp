package handler

import (
	"github.com/gorilla/sessions"
	"github.com/squiidz/stamp/module/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"time"
)

type Users struct {
	Username string
	Email    string
	Password string
	Create   time.Time
	Update   time.Time
	Friends  []string
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
	Position Location
	Picture  string
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

func RegisterHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		Templates.ExecuteTemplate(rw, "register.html", nil)
	} else if req.Method == "POST" {
		tempUser := Users{}
		newUser := Users{
			Username: req.FormValue("username"),
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
			Create:   time.Now(),
			Update:   time.Now(),
			Friends:  nil,
		}
		err = UCol.Find(bson.M{"username": newUser.Username}).One(&tempUser)
		if err != nil {
			UCol.Insert(&newUser)
			http.Redirect(rw, req, "/", http.StatusFound)
		} else {
			http.Redirect(rw, req, "/register", http.StatusFound)
		}
	}
}

func ProfilHandler(rw http.ResponseWriter, req *http.Request) {
	session, _ := Store.Get(req, "sessionCookie")
	data := session.Values["name"].(string)
	connUser := FindUser(data)

	Templates.ExecuteTemplate(rw, "profil.html", *connUser)
	logger.SimpleLog(req)
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	session, _ := Store.Get(req, "sessionCookie")
	data := session.Values["name"].(string)
	connUser := FindUser(data)

	Templates.ExecuteTemplate(rw, "index.html", *connUser)
	logger.SimpleLog(req)
}
