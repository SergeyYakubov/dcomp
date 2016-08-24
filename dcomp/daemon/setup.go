package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

var estimatorServer server.Server

var resources map[string]structs.Resource

func setServerConfiguration(srv *server.Server) error {
	srv.Host = "172.17.0.2"
	srv.Port = 27017
	return nil
}

func initializedb(name string) (db database.Agent, err error) {

	db = new(database.Mongodb)
	var srv server.Server
	if err := setServerConfiguration(&srv); err != nil {
		return nil, err
	}

	db.SetServer(&srv)
	db.SetDefaults("daemondbd")

	if err = db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func initialize() error {
	resources = make(map[string]structs.Resource)

	estimatorServer.Host = "localhost"
	estimatorServer.Port = 8002

	resources["Local"] = structs.Resource{Server: server.Server{"localhost", 8003}}

	return nil
}
