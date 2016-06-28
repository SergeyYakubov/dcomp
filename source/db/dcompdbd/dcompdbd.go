package main

import (
	"os"
	"stash.desy.de/scm/dc/version"
)

func main() {
	if ret := version.ShowVersion(os.Stdout, "dCompdbd"); ret {
		return
	}
	return
}
