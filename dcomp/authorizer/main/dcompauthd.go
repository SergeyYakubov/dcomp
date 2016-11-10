package main

import (
	"os"

	"github.com/sergeyyakubov/dcomp/dcomp/authorizer/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompauthd"); ret {
		return
	}
	daemon.Start()
}
