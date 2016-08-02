package main

import (
	"log"
	"net/http"
	"os"
	"stash.desy.de/scm/dc/main.git/dcomp/db/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"

	"stash.desy.de/scm/dc/main.git/dcomp/utils"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func initdb(name string) error {

	if err := database.Create(name); err != nil {
		return err
	}

	if err := database.SetServerConfiguration(); err != nil {
		return err
	}

	if err := database.Connect(); err != nil {
		return err
	}
	return nil
}

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompdbd"); ret {
		return
	}

	if err := initdb("mongodb"); err != nil {
		log.Fatal("mongodb: " + err.Error())
	}

	defer database.Close()
	mux := utils.NewRouter(daemon.ListRoutes)
	log.Fatal(http.ListenAndServe(":8001", mux))
}
