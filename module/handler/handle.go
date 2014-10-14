package handler

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

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
)

func init() {
	TemplatesLocation["index"] = "template/index.html"
	TemplatesLocation["watch"] = "template/watch.html"

	Templates.Delims("((", "))")
	for _, value := range TemplatesLocation {
		Templates.ParseFiles(value)
	}
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	Templates.ExecuteTemplate(rw, "index.html", "Naked")
	SimpleLog(req)
}

func WatchHandler(rw http.ResponseWriter, req *http.Request) {
	Templates.ExecuteTemplate(rw, "watch.html", "Here you are")
	SimpleLog(req)
}

func LocationHandler(rw http.ResponseWriter, req *http.Request) {
	m := Message{}
	loc := Location{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	err = c.Find(bson.M{"latitude": loc.Latitude, "longitude": loc.Longitude}).One(&m)
	if err != nil {
		log.Println("CANNOT FIND MESSAGE")
	}
	enco := json.NewEncoder(rw)
	enco.Encode(&m)
	c.Remove(bson.M{"latitude": loc.Latitude, "longitude": loc.Longitude})

}

func PlaceHandler(rw http.ResponseWriter, req *http.Request) {
	loc := Message{}
	data := json.NewDecoder(req.Body)
	data.Decode(&loc)

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Println("ERROR AT CONNECTING TO DB")
	}

	c := session.DB("message").C("new")
	c.Insert(&loc)

	/*log.Println("Friend : ", loc.Friends)
	log.Println("Message : ", loc.Message)
	log.Println("Latitude : ", loc.Latitude)
	log.Println("Longitude : ", loc.Longitude)*/
}
