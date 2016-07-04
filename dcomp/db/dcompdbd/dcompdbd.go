package main

import (
	"log"
	"net/http"
	"os"
	"stash.desy.de/scm/dc/main.git/dcomp/db/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func main() {
	if ret := version.ShowVersion(os.Stdout, "dCompdbd"); ret {
		return
	}
	daemon.SetServerConfiguration()
	mux := utils.NewRouter(daemon.ListRoutes)
	log.Fatal(http.ListenAndServe(":8001", mux))
}
