package daemon

import (
	"log"
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/resources"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
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

	log.Fatal(http.ListenAndServe(addr, server.ProcessHMACAuth(mux.ServeHTTP, key)))

	return nil
}
