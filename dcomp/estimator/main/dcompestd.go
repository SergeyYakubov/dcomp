package main

import (
	"os"

	"github.com/sergeyyakubov/dcomp/dcomp/estimator/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompestd"); ret {
		return
	}
	daemon.Start()
}
