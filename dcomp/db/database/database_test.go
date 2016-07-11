package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func TestCreate(t *testing.T) {
	db = new(mockdatabase)

	err := Create("mongodb")
	assert.NotNil(t, err)

	db = nil
	err = Create("notexist")
	assert.NotNil(t, err)
}

func TestConnect(t *testing.T) {
	db = new(mockdatabase)
	defer func() { db = nil }()
	dbServer.Host = "localhost"
	dbServer.Port = 27017
	err := Connect()
	assert.Nil(t, err)
}

func TestCreateRecord(t *testing.T) {
	db = new(mockdatabase)
	assert.NotNil(t, db, "should not be nil")

	s := struct{}{}

	id, err := CreateRecord(s)
	assert.Nil(t, err)
	assert.Equal(t, "578359205e935a20adb39a18", id, "record created")
	db = nil
	id, err = CreateRecord(s)
	assert.NotNil(t, err, "nil database")

}
func TestGetRecords(t *testing.T) {
	db = new(mockdatabase)
	defer func() { db = nil }()
	q := querryM{bson.ObjectIdHex("578359205e935a20adb39a18")}

	var records []structs.JobInfo
	err := GetRecords(&q, &records)
	assert.Nil(t, err, "Got records")
	assert.Equal(t, 1, len(records), "TestGetRecords length should be 1")
	assert.Equal(t, 1, records[0].Status, "TestGetRecords should return 1")

}
func TestGetAllRecords(t *testing.T) {
	db = new(mockdatabase)
	defer func() { db = nil }()

	var records []structs.JobInfo
	err := GetAllRecords(&records)
	assert.Nil(t, err, "Got all records")
	assert.Equal(t, 2, len(records), "TestGetAllRecords should return 1")
	assert.Equal(t, 1, records[0].Status, "TestGetRecords should return 1")
	assert.Equal(t, 1, records[1].Status, "TestGetRecords should return 2")
}
