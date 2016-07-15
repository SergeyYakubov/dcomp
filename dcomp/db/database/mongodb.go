package database

import (
	"time"

	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongodb struct {
	session *mgo.Session
	name    string
	col     string
	timeout time.Duration
}

func (db *mongodb) CreateRecord(s interface{}) (string, error) {

	if db.session == nil {
		return "", errors.New("database session not created")
	}
	c := db.session.DB(db.name).C(db.col)
	id := bson.NewObjectId()
	_, err := c.UpsertId(id, s)

	if err != nil {
		return "", err
	}
	err = c.UpdateId(id, bson.M{"$set": bson.M{"_hex_id": id.Hex()}})
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
	db.session = nil
}

func (db *mongodb) SetDefaults() {
	db.name = "daemondbd"
	db.col = "jobs"
	db.timeout = 10 * time.Second
}

func (db *mongodb) GetRecords(q interface{}, res interface{}) (err error) {

	c := db.session.DB(db.name).C(db.col)
	query := c.Find(q)

	n, err := query.Count()

	if err != nil {
		return err
	}

	if n == 0 {
		return nil
	}
	err = query.All(res)
	return err
}

func (db *mongodb) GetRecordByID(id string, res interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}
	return db.GetRecords(q, res)
}

func (db *mongodb) DeleteRecordByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}

	c := db.session.DB(db.name).C(db.col)
	_, err := c.RemoveAll(q)
	return err
}