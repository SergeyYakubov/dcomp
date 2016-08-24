package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/resources"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

func Start(res resources.Resource, db database.Agent, port string) error {

	p := resources.NewPlugin(res, db)
	if err := db.Connect(); err != nil {
		return err
	}
	defer db.Close()
	mux := utils.NewRouter(p.ListRoutes)
	log.Fatal(http.ListenAndServe(":"+port, mux))
	return nil
}
