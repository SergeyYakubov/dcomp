package jobdatabase

import (
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestCreateRecord(t *testing.T) {
	db := new(Mockdatabase)
	assert.NotNil(t, db, "should not be nil")

	s := struct{}{}

	id, err := db.CreateRecord("", s)
	assert.Nil(t, err)
	assert.Equal(t, "578359205e935a20adb39a18", id, "record created")

}

func TestPatchRecord(t *testing.T) {
	db := new(Mockdatabase)
	assert.NotNil(t, db, "should not be nil")

	s := struct{}{}

	err := db.PatchRecord("", s)
	assert.Nil(t, err)

	err = db.PatchRecord("give error", s)
	assert.NotNil(t, err)

}

func TestGetRecords(t *testing.T) {
	db := new(Mockdatabase)
	defer func() { db = nil }()
	q := querryM{bson.ObjectIdHex("578359205e935a20adb39a18")}

	var records []structs.JobInfo
	err := db.GetRecords(&q, &records)
	assert.Nil(t, err, "Got records")
	assert.Equal(t, 1, len(records), "TestGetRecords length should be 1")
	assert.Equal(t, 1, records[0].Status, "TestGetRecords should return 1")

}
func TestGetAllRecords(t *testing.T) {
	db := new(Mockdatabase)
	defer func() { db = nil }()

	var records []structs.JobInfo
	err := db.GetAllRecords(&records)
	assert.Nil(t, err, "Got all records")
	assert.Equal(t, 3, len(records), "TestGetAllRecords should return 3")
	assert.Equal(t, 1, records[0].Status, "TestGetRecords should return 1")
	assert.Equal(t, 1, records[1].Status, "TestGetRecords should return 1")
}

func TestDeleteRecord(t *testing.T) {
	db := new(Mockdatabase)
	defer func() { db = nil }()

	err := db.DeleteRecordByID("578359205e935a20adb39a18")
	assert.Nil(t, err, "Delete record")

	err = db.DeleteRecordByID("578359205e935a20adb39a19")
	assert.NotNil(t, err, "Not found record")
}
