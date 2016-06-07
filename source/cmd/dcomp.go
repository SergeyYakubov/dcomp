package main

import (
	"../cli"
	"../daemon"
	dcompversion "../version"
	"flag"
	"fmt"
)

var (
	flHelp    = flag.Bool("help", false, "Print usage")
	flVersion = flag.Bool("version", false, "Print version information")
)

func main() {

	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	if *flHelp || flag.NArg() == 0 {
		flag.Usage()
		return
	}

	if flag.Arg(0) == "daemon" {
		daemon.StartDaemon(flag.Args()[1:])
	} else {
		cli.Command(flag.Arg(0), flag.Args()[1:])
	}

}

func showVersion() {
	fmt.Printf("dComp version %s, build time %s\n", dcompversion.Version, dcompversion.BuildTime)
}
