package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/resources"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var resource resources.Resource

func Start(res resources.Resource, db database.Agent, addr string) error {
	resource = res
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resource.SetDb(db)
	mux := utils.NewRouter(listRoutes)
	log.Fatal(http.ListenAndServe(addr, mux))
	return nil
}
