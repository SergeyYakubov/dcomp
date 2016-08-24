// Package for central dComp daemon
package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var db database.Agent

func Start(args []string) {

	db, err := initializedb("mongodb")
	if err != nil {
		log.Fatal("cannot connect to mongodb: " + err.Error())
	}
	defer db.Close()

	initialize()

	mux := utils.NewRouter(listRoutes)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
