package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	configFileName = `daemon_test.yaml`
	c, _ := setConfiguration()

	c.BaseDir = "/aaa"
	err := c.check()
	assert.NotNil(t, err)

	c.BaseDir = "/bin/bash"
	err = c.check()
	assert.NotNil(t, err)

}

func TestSetConfiguration(t *testing.T) {

	configFileName = `main_test.yaml`

	c, err := setConfiguration()

	assert.Nil(t, err, "Should not be error")

	assert.Equal(t, c.Daemon.Addr, ":8006", "addr")
	assert.Equal(t, c.Daemon.Key, "12345", "key")

	assert.Equal(t, c.Database.Port, 27017, "port")

	configFileName = `aaa`

	c, err = setConfiguration()
	assert.NotNil(t, err, "Should be error")
}
