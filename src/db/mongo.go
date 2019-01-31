package db

import (
	"os"

	"github.com/globalsign/mgo"
)

func GetMongoSession() (*mgo.Session, error) {
	env := os.Getenv("ENVIRONMENT")
	var url string

	if env == "PRODUCTION" {
		url = "mongodb://api:dM6CYayNQu8qr9b@ds147003.mlab.com:47003/heroku_rvdsxf5j"
	} else {
		url = "localhost"
	}
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	session.SetSafe(&mgo.Safe{})
	return session, nil
}

func CloseMongoSession(session *mgo.Session) {
	session.Close()
}
