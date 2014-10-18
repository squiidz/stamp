package handler

import (
	"encoding/json"
	db "github.com/squiidz/stamp/module/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"time"
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

var (
	TemplatesLocation = map[string]string{}
	Templates         = template.New("main")
	//Connected         Users
)

func init() {
	TemplatesLocation["login"] = "template/login.html"
	TemplatesLocation["index"] = "template/index.html"
	TemplatesLocation["watch"] = "template/watch.html"

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
			cookie := http.Cookie{
				Name:  "stamp",
				Value: user.Username,
			}
			http.SetCookie(rw, &cookie)
			http.Redirect(rw, req, "/home", http.StatusFound)
		} else {
			http.Redirect(rw, req, "/", http.StatusFound)
		}
	}
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	value, _ := req.Cookie("stamp")
	data := value.Value
	Templates.ExecuteTemplate(rw, "index.html", data)
	SimpleLog(req)
}

func WatchHandler(rw http.ResponseWriter, req *http.Request) {
	Templates.ExecuteTemplate(rw, "watch.html", "Here you are")
	SimpleLog(req)
}

// Check User Position , and return Message if they exist for the current location
func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	m := []Message{}
	cookie, err := req.Cookie("stamp")
	if err != nil {
		log.Println("COOKIE DOESN'T EXIST")
	}
	username := cookie.Value
	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	err = c.Find(bson.M{"to.username": username}).All(&m)

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
			c.Remove(bson.M{"to.username": username, "message": mess.Message})
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

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	c.Insert(&message)

}

func CheckUser(user Users) bool {
	session, err := mgo.Dial(db.ServerAddr)
	defer session.Close()
	if err != nil {
		log.Println("Cannot connect to Database")
	}
	collection := session.DB("account").C("user")
	err = collection.Find(bson.M{"user": user.Username, "password": user.Password}).One(&user)
	if err != nil {
		return false
	} else {
		return true
	}
}

func PositionValid(message *Message, location *Location) bool {
	var zone float64 = 0.0000200

	if (location.Latitude-message.Position.Latitude) < zone || (location.Latitude-message.Position.Latitude) > (zone-zone*2) && (location.Latitude-message.Position.Latitude) < zone {
		if (location.Longitude-message.Position.Longitude) < zone || (location.Longitude-message.Position.Longitude) > (zone-zone*2) && (location.Longitude-message.Position.Longitude) < zone {
			/*log.Println("TRUE")
			fmt.Printf("DIFF LAT: %.6f\n", (location.Latitude - message.Latitude))
			fmt.Printf("DIFF LONG: %.6f\n", (location.Longitude - message.Longitude))*/
			return true
		} else {
			/*fmt.Printf("DIFF : %.6f\n", (location.Longitude - message.Longitude))
			log.Println("LONGITUDE FALSE")*/
			return false
		}
	} else {
		/*fmt.Printf("DIFF : %.6f\n", (location.Latitude - message.Latitude))
		log.Println("LATITUDE FALSE")*/
		return false
	}
}

func ForgeCookie(user *Users) *http.Cookie {
	delay := time.Now()

	newCookie := &http.Cookie{
		Name:    "connected",
		Value:   user.Username,
		Expires: delay,
	}
	return newCookie
}
