package handler

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

func CheckUser(user Users) bool {
	if err != nil {
		log.Println("Cannot connect to Database")
	}
	err = UCol.Find(bson.M{"user": user.Username, "password": user.Password}).One(&user)
	if err != nil {
		return false
	} else {
		return true
	}
}

func PositionValid(message *Message, location *Location) bool {
	var zone float64 = 1.0005200

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
