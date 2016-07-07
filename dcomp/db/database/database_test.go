package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	db = new(mockdatabase)

	err := Create("mongodb")
	assert.NotNil(t, err)

	db = nil
	err = Create("notexist")
	assert.NotNil(t, err)
}

func TestCreateRecord(t *testing.T) {
	db = new(mockdatabase)
	assert.NotNil(t, db, "should not be nil")

	s := struct{}{}

	id, err := CreateRecord(s)
	assert.Nil(t, err)
	assert.Equal(t, "1", id, "record created")
	db = nil
	id, err = CreateRecord(s)
	assert.NotNil(t, err, "nil database")

}

func TestConnect(t *testing.T) {
	db = new(mockdatabase)
	dbServer.Host = "localhost"
	dbServer.Port = 27017
	err := Connect()
	assert.Nil(t, err)
	db = nil
}
