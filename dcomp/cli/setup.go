package cli

import (
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

var Server server.Server

func SetServerConfiguration() error {
	Server.Host = "localhost"
	Server.Port = 8000
	return nil
}
