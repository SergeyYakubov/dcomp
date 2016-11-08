package daemon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetConfiguration(t *testing.T) {
	setConfiguration()
	assert.Equal(t, "172.17.0.2", dbServer.Host, "Host")
	assert.Equal(t, 27017, dbServer.Port, "Port")
	assert.Equal(t, "db3415", dbServer.Key, "Key")
	assert.Equal(t, ":8001", addr, "addr")
	assert.Equal(t, "localhost", resources["Local"].Server.Host, "local res host")
	assert.Equal(t, 8006, resources["Local"].Server.Port, "local res port")
	assert.Equal(t, "lplgin12", resources["Local"].Server.Key, "key")
	assert.Equal(t, "localhost", estimatorServer.Host, "est host")
	assert.Equal(t, 8005, estimatorServer.Port, "est port")
	assert.Equal(t, "estim003", estimatorServer.Key, "est key")

}
