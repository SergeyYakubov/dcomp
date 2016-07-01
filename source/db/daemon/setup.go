package daemon

import (
	"stash.desy.de/scm/dc/server"
)

var Server server.Srv

func SetServerConfiguration() error {
	Server.Host = "localhost"
	Server.Port = 8002
	return nil
}