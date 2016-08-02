package cli

import (
	"stash.desy.de/scm/dc/main.git/dcomp/server"
)

var daemon server.Server

// SetDaemonConfiguration reads configuration file with daemon location
func SetDaemonConfiguration() error {
	daemon.Host = "localhost"
	daemon.Port = 8000
	return nil
}
