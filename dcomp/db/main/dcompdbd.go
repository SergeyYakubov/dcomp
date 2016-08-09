package main

import (
	"log"
	"os"

	"stash.desy.de/scm/dc/main.git/dcomp/db/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"

	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func setServerConfiguration(srv *server.Server) error {
	srv.Host = "172.17.0.2"
	srv.Port = 27017
	return nil
}

func initdb(name string) (db database.Agent, err error) {

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

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompdbd"); ret {
		return
	}
	var db database.Agent
	db, err := initdb("mongodb")
	if err != nil {
		log.Fatal("mongodb: " + err.Error())
	}
	defer db.Close()

	daemon.Start(db)
}
