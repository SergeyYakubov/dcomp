package database

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

//these tests assume that mongodb server is running on 172.17.0.2:27017 (best to use Docker container)

func TestMdbConnect(t *testing.T) {
	err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb should be set")

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	err = Connect()
	assert.Nil(t, err)
	db.Close()

	db.(*mongodb).timeout = time.Second / 10
	dbServer.Host = ""
	dbServer.Port = 27017
	err = Connect()
	assert.NotNil(t, err)
	db = nil
}

func TestMdbCreateRecord(t *testing.T) {

	err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	s := structs.JobInfo{structs.JobDescription{}, "dummyid", 1}

	err = Connect()
	assert.Nil(t, err, "connected to database")

	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	id, err = db.CreateRecord(nil)
	assert.NotNil(t, err, "nil record")
	db.Close()

	_, err = db.CreateRecord(&s)
	assert.NotNil(t, err, "closed database")

	db = nil
}

func TestMdbGetRecords(t *testing.T) {
	err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	err = Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{structs.JobDescription{"name", "script", 20}, "dummyid", 1}
	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id, "normal record")

	var records []structs.JobInfo

	q := bson.M{"jobdescription.imagename": "name", "_id": bson.ObjectIdHex(id)}

	err = db.GetRecords(q, &records)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(records), "TestMdbGetRecords should return 1")
	assert.Equal(t, id, records[0].Id, "TestMdbGetRecords should return same id")
	db = nil

}

func TestMdbGetRecordByID(t *testing.T) {
	err := Create("mongodb")
	assert.Nil(t, err, "Pointer to mongodb shoud be set")

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	err = Connect()
	assert.Nil(t, err, "connected to database")

	s := structs.JobInfo{structs.JobDescription{}, "dummyid", 1}
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

	db = nil

}
