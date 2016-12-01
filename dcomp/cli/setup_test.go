package cli

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetServerConfiguration(t *testing.T) {
	SetDaemonConfiguration()
	assert.Equal(t, "172.18.0.3", daemon.Host, "")
	assert.Equal(t, 8001, daemon.Port, "")

	_, ok := daemon.GetAuth().(*server.BasicAuth)
	assert.Equal(t, ok, true, "auth set")

}
