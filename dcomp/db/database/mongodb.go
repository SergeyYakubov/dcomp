package database

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongodb struct {
	session *mgo.Session
	name    string
	col     string
	timeout time.Duration
}

type Person struct {
	Name  string
	Phone string
}

func (db *mongodb) CreateRecord(s interface{}) (string, error) {

	c := db.session.DB(db.name).C(db.col)
	id := bson.NewObjectId()
	_, err := c.UpsertId(id, s)

	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (db *mongodb) Connect(url string) error {
	var err error
	db.session, err = mgo.DialWithTimeout(url, db.timeout)
	return err
}

func (db *mongodb) Close() {
	db.session.Close()
}

func (db *mongodb) SetDefaults() {
	db.name = "daemondbd"
	db.col = "jobs"
	db.timeout = 10 * time.Second
}
