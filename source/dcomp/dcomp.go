package main

import (
	"flag"
	"fmt"
	"os"
	"stash.desy.de/scm/dc/cli"
	"stash.desy.de/scm/dc/daemon"
	dcompversion "stash.desy.de/scm/dc/version"
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
		fmt.Fprintln(cli.OutBuf, "\nCommands:")
		cli.PrintAllCommands()
		return
	}

	if flag.Arg(0) == "daemon" {
		daemon.StartDaemon(flag.Args()[1:])
	} else {
		if err := cli.SetServerConfiguration(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := cli.Command(flag.Arg(0), flag.Args()[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func showVersion() {
	fmt.Fprintf(cli.OutBuf, "dComp version %s, build time %s\n", dcompversion.Version, dcompversion.BuildTime)
}
