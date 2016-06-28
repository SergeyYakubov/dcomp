package cli

import (
	"stash.desy.de/scm/dc/server"
)

var Server server.Srv

func SetServerConfiguration() error {
	Server.Host = "localhost"
	Server.Port = 8000
	return nil
}
