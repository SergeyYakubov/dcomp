package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetConfiguration(t *testing.T) {

	tconf := configFile
	configFile = "config_test.yaml"
	setDaemonConfiguration()
	assert.Equal(t, ":8007", c.Daemon.Addr, "port")
	assert.Equal(t, "/etc/dcomp/cert/certauth.pem", c.Daemon.Certfile, "certfile")
	assert.Equal(t, "/etc/dcomp/cert/keyauth.pem", c.Daemon.Keyfile, "certfile")
	assert.Equal(t, "Basic", c.Authorization[0], "Allowed auth")
	assert.Equal(t, "Negotiate", c.Authorization[1], "Allowed auth")

	configFile = tconf
	setDaemonConfiguration()

}
