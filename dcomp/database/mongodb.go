package database

import (
	"time"

	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

type Mongodb struct {
	session *mgo.Session
	name    string
	col     string
	timeout time.Duration
	srv     *server.Server
}

// CreateRecord changes a database record with given id. s should be is an object that
// mgo understands (go struct is OK)
func (db *Mongodb) PatchRecord(id string, s interface{}) error {
	if err := checkID(id); err != nil {
		return err
	}

	c := db.session.DB(db.name).C(db.col)
	return c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": s})
}

// CreateRecord creates a database record with new unique id. s should be is an object that
// mgo understands (go struct is OK)
func (db *Mongodb) CreateRecord(given_id string, s interface{}) (string, error) {

	if db.session == nil {
		return "", errors.New("database session not created")
	}
	c := db.session.DB(db.name).C(db.col)

	var id bson.ObjectId
	if given_id == "" {
		// create new unique id
		id = bson.NewObjectId()
	} else {
		if bson.IsObjectIdHex(given_id) {
			id = bson.ObjectIdHex(given_id)
		} else {
			return "", errors.New("Bad id format")
		}
	}

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

func (db *Mongodb) SetServer(srv *server.Server) {
	db.srv = srv
}

func (db *Mongodb) Connect() error {
	var err error
	db.session, err = mgo.DialWithTimeout(db.srv.FullName(), db.timeout)
	return err
}

func (db *Mongodb) Close() {
	db.session.Close()
	db.session = nil
}

func (db *Mongodb) SetDefaults(name ...interface{}) {
	if len(name) > 0 {
		db.name = name[0].(string)
	}
	db.col = "jobs"
	db.timeout = 10 * time.Second
}

// GetRecords issues a request to mongodb. q should be a bson.M object or go struct with fields to match
// returns
func (db *Mongodb) GetRecords(q interface{}, res interface{}) (err error) {

	c := db.session.DB(db.name).C(db.col)
	query := c.Find(q)

	n, err := query.Count()

	if err != nil && err != io.EOF {
		return err
	}

	if n == 0 {
		return nil
	}
	err = query.All(res)
	return err
}

// GetAllRecords returns all records
func (db *Mongodb) GetAllRecords(res interface{}) (err error) {
	return db.GetRecords(nil, res)
}

func checkID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}
	return nil
}

func (db *Mongodb) GetRecordByID(id string, res interface{}) error {
	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}
	return db.GetRecords(q, res)
}

func (db *Mongodb) DeleteRecordByID(id string) error {
	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}

	c := db.session.DB(db.name).C(db.col)
	_, err := c.RemoveAll(q)
	return err
}
