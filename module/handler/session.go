package handler

import (
	"crypto/md5"
	"github.com/garyburd/redigo/redis"
	"log"
)

var (
	c   redis.Conn
	err error
)

func init() {
	c, err = redis.Dial("tcp", "192.168.0.104:6379")
	if err != nil {
		log.Println(err)
	}
}

func CheckSession(md string) string {
	name, err := redis.String(c.Do("GET", "id:"+md))
	if err != nil {
		log.Println("[*] ID NOT CREATED")
		return ""
	}
	return name
}

func MakeSession(user *Users) string {
	cypho := md5.Sum([]byte(user.Username + user.Password))
	_, err = c.Do("Set", "id:"+string(cypho[:16]), user.Username)
	result := string(cypho[:15])
	return result
}

/*
func GetUsername(md string) string {
	if err != nil {
		log.Println(err)
	}
	err = c.Do("Get", "id:"+md)

}
*/
