package handler

import (
	//"github.com/gorilla/sessions"
	"github.com/squiidz/stamp/module/logger"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"html/template"
	"net/http"
	"time"
)

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
	connUser := &Users{
		Username: CookieValue(req, "sessionCookie"),
	}
	err := connUser.FindUser()
	if err != nil {
		http.Redirect(rw, req, "/", http.StatusFound)
	}

	Templates.ExecuteTemplate(rw, "profil.html", *connUser)
	logger.SimpleLog(req)
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	connUser := &Users{
		Username: CookieValue(req, "sessionCookie"),
	}
	err := connUser.FindUser()
	if err != nil {
		http.Redirect(rw, req, "/", http.StatusFound)
	}

	Templates.ExecuteTemplate(rw, "index.html", *connUser)
	logger.SimpleLog(req)
}
