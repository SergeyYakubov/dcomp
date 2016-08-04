package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

var dbServer server.Server
var estimatorServer server.Server

func setServers() error {
	dbServer.Host = "localhost"
	dbServer.Port = 8001
	estimatorServer.Host = "localhost"
	estimatorServer.Port = 8002
	return nil
}
