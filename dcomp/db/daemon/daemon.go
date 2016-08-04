// Package for database interface daemon
package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

func Start() {
	mux := utils.NewRouter(listRoutes)
	log.Fatal(http.ListenAndServe(":8001", mux))
}
