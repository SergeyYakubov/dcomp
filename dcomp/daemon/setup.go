package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

var dbServer server.Server
var estimatorServer server.Server

var resources map[string]structs.Resource

func initialize() error {
	resources = make(map[string]structs.Resource)
	dbServer.Host = "localhost"
	dbServer.Port = 8001
	estimatorServer.Host = "localhost"
	estimatorServer.Port = 8002

	resources["Local"] = structs.Resource{Server: server.Server{"localhost", 8003}}

	return nil
}
