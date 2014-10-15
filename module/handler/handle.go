package handler

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	//"fmt"
	"log"
	"net/http"
	//"math"
)

type Users struct {
	Username string
	Password string
}

type Location struct {
	Longitude float64
	Latitude  float64
}

type Message struct {
	Friends   string
	Message   string
	Longitude float64
	Latitude  float64
}

var (
	TemplatesLocation = map[string]string{}
	Templates         = template.New("main")
	Connected string
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
	} else if req.Method == "POST"{
		user := Users{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}
		if CheckUser(user) {
			Connected = user.Username
			http.Redirect(rw, req, "/home", http.StatusFound)
		}
	}
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	Templates.ExecuteTemplate(rw, "index.html", nil)
	SimpleLog(req)
}

func WatchHandler(rw http.ResponseWriter, req *http.Request) {
	Templates.ExecuteTemplate(rw, "watch.html", "Here you are")
	SimpleLog(req)
}

// Check User Position , and return Message if they exist for the current location
func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	m := []Message{}

	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)
	log.Println(loc.Latitude)
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	err = c.Find(bson.M{"friends": Connected}).All(&m)
	if err != nil {
		log.Println("CANNOT FIND MESSAGE")
		log.Println(err) // Return nothing if no messages
	}
	// Encode to Json messages found
	enco := json.NewEncoder(rw)
	for _, mess := range m {
		if PositionValid(&mess, &loc){
			//valMes = append(valMes, mess)
			log.Println(mess.Message)
			enco.Encode(&mess)
			c.Remove(bson.M{"friends": Connected, "message": mess.Message})
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

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	c.Insert(&message)

	/*log.Println("Friend : ", loc.Friends)
	log.Println("Message : ", loc.Message)
	log.Println("Latitude : ", loc.Latitude)
	log.Println("Longitude : ", loc.Longitude)*/
}

func CheckUser(user Users) bool {
	return true
}

func PositionValid(message *Message, location *Location) bool{
	var zone float64 = 0.0000200
	//log.Println(zone)
	/*location.Latitude = math.Abs(location.Latitude)
	location.Longitude = math.Abs(location.Longitude)
	message.Latitude = math.Abs(message.Latitude)
	message.Longitude = math.Abs(message.Longitude)*/

	log.Println("Check Validity")
	if (location.Latitude - message.Latitude) < zone || (location.Latitude - message.Latitude) > (zone - zone*2) && (location.Latitude - message.Latitude) < zone{
		if (location.Longitude - message.Longitude) < zone || (location.Longitude - message.Longitude) > (zone - zone*2) && (location.Longitude - message.Longitude) < zone {
			/*log.Println("TRUE")
			fmt.Printf("DIFF LAT: %.6f\n", (location.Latitude - message.Latitude))
			fmt.Printf("DIFF LONG: %.6f\n", (location.Longitude - message.Longitude))*/
			return true
		}else {
			/*fmt.Printf("DIFF : %.6f\n", (location.Longitude - message.Longitude))
			log.Println("LONGITUDE FALSE")*/
			return false
		}
	}else {
		/*fmt.Printf("DIFF : %.6f\n", (location.Latitude - message.Latitude))
		log.Println("LATITUDE FALSE")*/
		return false
	}
}
