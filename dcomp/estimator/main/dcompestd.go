package main

import (
	"log"
	"net/http"
	"os"

	"stash.desy.de/scm/dc/main.git/dcomp/estimator/daemon"

	"stash.desy.de/scm/dc/main.git/dcomp/utils"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompestd"); ret {
		return
	}

	mux := utils.NewRouter(daemon.ListRoutes)
	log.Fatal(http.ListenAndServe(":8002", mux))
}
