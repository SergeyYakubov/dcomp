package main

import (
	dcompversion "../version"
	"fmt"
	flag "github.com/docker/docker/pkg/mflag"
)

var (
	flHelp    = flag.Bool([]string{"h", "-help"}, false, "Print usage")
	flVersion = flag.Bool([]string{"v", "-version"}, false, "Print version information and quit")
)

func main() {

	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	if *flHelp {
		flag.Usage()
		return
	}
}

func Dummy() int {
	return 2
}

func showVersion() {
	fmt.Printf("dComp version %s, build at %s\n", dcompversion.Version, dcompversion.BuildTime)
}
