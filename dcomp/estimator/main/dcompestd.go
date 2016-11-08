package main

import (
	"os"

	"github.com/dcomp/dcomp/estimator/daemon"
	"github.com/dcomp/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompestd"); ret {
		return
	}
	daemon.Start()
}
