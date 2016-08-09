package main

import (
	"log"
	"os"

	"stash.desy.de/scm/dc/main.git/dcomp/db/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"

	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func initdb(name string) (db database.Agent, err error) {
	if db, err = database.Create(name); err != nil {
		return nil, err
	}

	if err = db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompdbd"); ret {
		return
	}
	var db database.Agent
	db, err := initdb("mongodb")
	if err != nil {
		log.Fatal("mongodb: " + err.Error())
	}
	defer db.Close()

	daemon.Start(db)
}
