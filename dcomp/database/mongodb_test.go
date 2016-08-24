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

func initdb() *Mongodb {
	db := new(Mongodb)

	var dbServer server.Server

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	db.SetServer(&dbServer)
	db.SetDefaults("daemondbd")

	return db
}

func TestMdbConnect(t *testing.T) {

	db := initdb()
	err := db.Connect()

	assert.Nil(t, err)
	db.Close()

	db.timeout = time.Second / 10
	db.srv.Host = ""
	db.srv.Port = 27017
	err = db.Connect()
	assert.NotNil(t, err)
}

func TestMdbCreateRecord(t *testing.T) {

	db := initdb()

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "dummyid", Status: 1}

	err := db.Connect()
	assert.Nil(t, err, "connected to database")

	id, err := db.CreateRecord("", &s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	id, err = db.CreateRecord("", nil)
	assert.NotNil(t, err, "nil record")
	db.Close()

	_, err = db.CreateRecord("", &s)
	assert.NotNil(t, err, "closed database")

}

func TestMdbGetRecords(t *testing.T) {

	db := initdb()

	err := db.Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{JobDescription: structs.JobDescription{ImageName: "name", Script: "script", NCPUs: 20}, Id: "dummyid", Status: 1}
	id, err := db.CreateRecord("", &s)
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

	db := initdb()
	err := db.Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "dummyid", Status: 1}
	id, err := db.CreateRecord("", &s)
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

func TestMdbPatchRecord(t *testing.T) {

	db := initdb()
	err := db.Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "dummyid", Status: 0}
	id, err := db.CreateRecord("", &s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	s.Status = 2
	s.Resource = "hello"
	err = db.PatchRecord(id, &s)
	assert.Nil(t, err)

	var records []structs.JobInfo

	err = db.GetRecordByID(id, &records)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(records), "TestMdbPatchRecord should return 1")
	assert.Equal(t, 2, records[0].Status, "TestMdbPatchRecord should return status 2")
	assert.Equal(t, "hello", records[0].Resource, "TestMdbPatchRecord should return resource hello")

	err = db.PatchRecord("aaa", &s)
	assert.NotNil(t, err)

	err = db.PatchRecord(id, nil)
	assert.NotNil(t, err)

	db.Close()
}