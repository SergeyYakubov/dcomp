package database

import (
	"time"

	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

type mongodb struct {
	session *mgo.Session
	name    string
	col     string
	timeout time.Duration
	srv     *server.Server
}

// CreateRecord creates a database record with new unique id. s should be is an object that
// mgo understands (go struct is OK)
func (db *mongodb) CreateRecord(s interface{}) (string, error) {

	if db.session == nil {
		return "", errors.New("database session not created")
	}
	c := db.session.DB(db.name).C(db.col)

	// create new unique id
	id := bson.NewObjectId()

	_, err := c.UpsertId(id, s)

	if err != nil {
		return "", err
	}
	// we keep both object id for faster search and its hex representation which can be passed to clients
	// within JSON struct
	err = c.UpdateId(id, bson.M{"$set": bson.M{"_hex_id": id.Hex()}})
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (db *mongodb) SetServer(srv *server.Server) {
	db.srv = srv
}

func (db *mongodb) Connect() error {
	var err error
	db.session, err = mgo.DialWithTimeout(db.srv.FullName(), db.timeout)
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

// GetRecords issues a request to mongodb. q should be a bson.M object or go struct with fields to match
// returns
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

// GetAllRecords returns all records
func (db *mongodb) GetAllRecords(res interface{}) (err error) {
	return db.GetRecords(nil, res)
}

func checkID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}
	return nil
}

func (db *mongodb) GetRecordByID(id string, res interface{}) error {
	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}
	return db.GetRecords(q, res)
}

func (db *mongodb) DeleteRecordByID(id string) error {
	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}

	c := db.session.DB(db.name).C(db.col)
	_, err := c.RemoveAll(q)
	return err
}
