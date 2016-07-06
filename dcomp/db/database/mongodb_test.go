package database

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
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

	s := struct {
		name   string
		number int
	}{"hello", 1}

	err = Connect()
	assert.Nil(t, err)

	id, err := db.CreateRecord(&s)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)

	_, err = db.CreateRecord("aaa")
	assert.NotNil(t, err)

	db.Close()
	db = nil
}
