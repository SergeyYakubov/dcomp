package daemon

import (
	"log"
	"net/http"

	"github.com/dcomp/dcomp/database"
	"github.com/dcomp/dcomp/resources"
	"github.com/dcomp/dcomp/utils"
)

var resource resources.Resource

func Start(res resources.Resource, db database.Agent, addr, key string) error {
	resource = res
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resource.SetDb(db)
	mux := utils.NewRouter(listRoutes)

	log.Fatal(http.ListenAndServe(addr, utils.Auth(mux.ServeHTTP, key)))

	return nil
}
