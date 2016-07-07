package database

import (
	"errors"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

type database interface {
	CreateRecord(s interface{}) (string, error)
	Connect(string) error
	SetDefaults()
	Close()
}

var db database

var dbServer server.Srv

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
}

func CreateRecord(s interface{}) (string, error) {
	if db == nil {
		return "", errors.New("database not set")
	}

	return db.CreateRecord(s)
}
