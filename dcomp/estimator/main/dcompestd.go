package main

import (
	"os"

	"stash.desy.de/scm/dc/main.git/dcomp/estimator/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompestd"); ret {
		return
	}

	daemon.Start()
}
