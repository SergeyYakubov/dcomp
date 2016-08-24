// Package containes a database objects and functions to work with it.
// db is an interface to a specific implementation (currently implemented mongodb and mockdatabase used for tests)
package database

import "stash.desy.de/scm/dc/main.git/dcomp/server"

type Agent interface {
	CreateRecord(string, interface{}) (string, error)
	GetAllRecords(interface{}) error
	GetRecords(interface{}, interface{}) error
	GetRecordByID(string, interface{}) error
	DeleteRecordByID(string) error
	Connect() error
	SetDefaults(name ...interface{})
	SetServer(*server.Server)
	Close()
}
