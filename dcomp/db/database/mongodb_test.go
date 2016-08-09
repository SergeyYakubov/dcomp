package database

import (
	"testing"

	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

//these tests assume that mongodb server is running on 172.17.0.2:27017 (best to use Docker container)

func TestMdbConnect(t *testing.T) {
	db, err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb should be set")

	var dbServer server.Server

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	err = db.Connect()
	db.SetServer(&dbServer)
	assert.Nil(t, err)
	db.Close()

	db.(*mongodb).timeout = time.Second / 10
	dbServer.Host = ""
	dbServer.Port = 27017
	db.SetServer(&dbServer)
	err = db.Connect()
	assert.NotNil(t, err)
}

func TestMdbCreateRecord(t *testing.T) {

	db, err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")

	var dbServer server.Server

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017
	db.SetServer(&dbServer)

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "dummyid", Status: 1}

	err = db.Connect()
	assert.Nil(t, err, "connected to database")

	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	id, err = db.CreateRecord(nil)
	assert.NotNil(t, err, "nil record")
	db.Close()

	_, err = db.CreateRecord(&s)
	assert.NotNil(t, err, "closed database")

}

func TestMdbGetRecords(t *testing.T) {
	db, err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")
	var dbServer server.Server

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017
	db.SetServer(&dbServer)

	err = db.Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{JobDescription: structs.JobDescription{"name", "script", 20}, Id: "dummyid", Status: 1}
	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	var records []structs.JobInfo

	q := bson.M{"jobdescription.imagename": "name", "_id": bson.ObjectIdHex(id)}

	err = db.GetRecords(q, &records)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(records), "TestMdbGetRecords should return 1")
	assert.Equal(t, id, records[0].Id, "TestMdbGetRecords should return same id")

	err = db.GetRecords("aaa", &records)
	assert.NotNil(t, err, "wrong querry")
	db.Close()

}

func TestMdbGetRecordByID(t *testing.T) {
	db, err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")

	var dbServer server.Server
	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017
	db.SetServer(&dbServer)

	err = db.Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "dummyid", Status: 1}
	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	var records []structs.JobInfo

	err = db.GetRecordByID(id, &records)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(records), "TestMdbGetRecords should return 1")
	assert.Equal(t, id, records[0].Id, "TestMdbGetRecords should return same id")

	err = db.GetRecordByID("aaa", &records)
	assert.NotNil(t, err)

	db.Close()
}
