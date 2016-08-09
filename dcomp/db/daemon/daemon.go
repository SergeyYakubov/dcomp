// Package for database interface daemon
package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var db database.Agent

func Start(database database.Agent) {
	db = database
	mux := utils.NewRouter(listRoutes)
	log.Fatal(http.ListenAndServe(":8001", mux))
}
