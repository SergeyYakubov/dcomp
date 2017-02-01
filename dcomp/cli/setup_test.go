package cli

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetServerConfiguration(t *testing.T) {
	SetDaemonConfiguration()
	assert.Equal(t, "localhost", daemon.Host, "")
	assert.Equal(t, 8001, daemon.Port, "")

	_, ok := daemon.GetAuth().(*server.GSSAPIAuth)
	assert.Equal(t, ok, true, "auth set")

}
