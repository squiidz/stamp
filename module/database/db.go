package database

import (
	"gopkg.in/mgo.v2"
)

func ConnectDB(server string) (*mgo.Session, error){
   if session, err = mgo.Dial(server); if err != nil {
        return err
   }
   return *session
}
