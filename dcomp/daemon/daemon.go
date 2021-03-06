// Package for central dComp daemon
package daemon

import (
	"log"
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

var db jobdatabase.Agent

func Start(args []string) {
	var err error

	err = setConfiguration()
	if err != nil {
		log.Fatal("cannot setup dcompd: " + err.Error())
	}

	db, err = connectDb("mongodb")
	if err != nil {
		log.Fatal("cannot connect to mongodb: " + err.Error())
	}
	defer db.Close()

	mux := utils.NewRouter(listRoutes)
	log.Fatal(http.ListenAndServeTLS(settings.Addr, settings.Certfile, settings.Keyfile,
		ProcessUserAuth(mux.ServeHTTP)))
}
