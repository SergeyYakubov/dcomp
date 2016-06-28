package daemon

import (
	"stash.desy.de/scm/dc/server"
)

var DBServer server.Srv

func SetServerConfiguration() error {
	DBServer.Host = "localhost"
	DBServer.Port = 8001
	return nil
}
