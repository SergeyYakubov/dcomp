// Package containes a database objects and functions to work with it.
// db is an interface to a specific implementation (currently implemented mongodb and mockdatabase used for tests)
package database

import (
	"errors"

	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

type Agent interface {
	CreateRecord(interface{}) (string, error)
	GetAllRecords(interface{}) error
	GetRecords(interface{}, interface{}) error
	GetRecordByID(string, interface{}) error
	DeleteRecordByID(string) error
	Connect() error
	SetDefaults()
	SetServer(*server.Server)
	Close()
}

func setServerConfiguration(srv *server.Server) error {
	srv.Host = "172.17.0.2"
	srv.Port = 27017
	return nil
}

func Create(name string) (Agent, error) {
	var db Agent
	switch name {
	case "mongodb":
		db = new(mongodb)
	default:
		return nil, errors.New("database " + name + " not found")
	}

	var srv server.Server
	if err := setServerConfiguration(&srv); err != nil {
		return nil, err
	}
	db.SetServer(&srv)
	db.SetDefaults()
	return db, nil

}
