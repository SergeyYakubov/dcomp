package daemon

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetConfiguration(t *testing.T) {

	tconf := configFile
	configFile = "setup_test.yaml"
	setConfiguration()
	assert.Equal(t, "172.17.0.2", dbServer.Host, "Host")
	assert.Equal(t, 27017, dbServer.Port, "Port")
	assert.Equal(t, "db3415", dbServer.GetAuth().(*server.HMACAuth).Key, "Key")
	assert.Equal(t, ":8001", settings.Addr, "addr")
	assert.Equal(t, "/etc/dcomp/cert/certd.pem", settings.Certfile, "certfile")
	assert.Equal(t, "/etc/dcomp/cert/keyd.pem", settings.Keyfile, "certfile")
	s := resources["local"].Server
	dm := resources["local"].DataManager
	assert.Equal(t, "localhost", s.Host, "local res host")
	assert.Equal(t, 8006, s.Port, "local res port")

	assert.Equal(t, "localhost", dm.Host, "local res host")
	assert.Equal(t, 8008, dm.Port, "local res port")
	assert.Equal(t, "dcompdmd45", dm.GetAuth().(*server.JWTAuth).Key, "key")

	assert.Equal(t, "lplgin12", s.GetAuth().(*server.HMACAuth).Key, "key")
	assert.Equal(t, "localhost", estimatorServer.Host, "est host")
	assert.Equal(t, 8005, estimatorServer.Port, "est port")
	assert.Equal(t, "estim003", estimatorServer.GetAuth().(*server.HMACAuth).Key, "est key")
	assert.Equal(t, "localhost", authServer.Host, "auth host")
	assert.Equal(t, 8007, authServer.Port, "auth port")
	assert.Equal(t, "auth14", authServer.GetAuth().(*server.HMACAuth).Key, "auth key")

	configFile = tconf
	setConfiguration()

}
