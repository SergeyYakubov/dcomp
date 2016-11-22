package main

import (
	"flag"
	"log"
	"os"

	"github.com/sergeyyakubov/dcomp/dcomp/datamanager/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

func main() {

	if ret := version.ShowVersion(os.Stdout, "dcompdmd"); ret {
		return
	}

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("configuration file not set. Usage: dcompdmd <config file name>")
	}

	daemon.Start(flag.Arg(0))
}
