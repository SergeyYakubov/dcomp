package main

import (
	"log"

	"os"

	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/resources/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/resources/local"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dComp local resource plugin"); ret {
		return
	}

	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugin")
	port := "8003"

	var res = new(local.Resource)

	log.Fatal(daemon.Start(res, db, port))
	return
}
