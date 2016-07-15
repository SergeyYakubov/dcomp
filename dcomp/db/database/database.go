package database

import (
	"errors"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

type agent interface {
	CreateRecord(interface{}) (string, error)
	GetRecords(interface{}, interface{}) error
	GetRecordByID(string, interface{}) error
	DeleteRecordByID(string) error
	Connect(string) error
	SetDefaults()
	Close()
}

var db agent

var dbServer server.Server

func SetServerConfiguration() error {
	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017
	return nil
}

func Create(name string) error {
	if db != nil {
		return errors.New("database already created")
	}
	switch name {
	case "mongodb":
		db = new(mongodb)
	default:
		return errors.New("database " + name + " not found")
	}

	db.SetDefaults()
	return nil
}

func Connect() error {
	return db.Connect(dbServer.HostPort())
}

func Close() {
	db.Close()
	db = nil
}

func CreateRecord(s interface{}) (string, error) {
	if db == nil {
		return "", errors.New("database not set")
	}

	return db.CreateRecord(s)
}

func GetRecords(q interface{}, res interface{}) (err error) {
	if db == nil {
		return errors.New("database not set")
	}
	return db.GetRecords(q, res)
}

func GetAllRecords(res interface{}) (err error) {
	return GetRecords(nil, res)
}

func GetRecordById(id string, res interface{}) (err error) {
	if db == nil {
		return errors.New("database not set")
	}
	return db.GetRecordByID(id, res)
}

func DeleteRecordById(id string) (err error) {
	if db == nil {
		return errors.New("database not set")
	}
	return db.DeleteRecordByID(id)
}