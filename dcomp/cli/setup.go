package cli

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

var daemon server.Server

type config struct {
	Dcompd struct {
		Port int
		Host string
	}
}

// SetDaemonConfiguration reads configuration file with daemon location
func SetDaemonConfiguration() error {
	fname := `/etc/dcomp/dcomp.yaml`

	var c config

	err := utils.ReadYaml(fname, &c)

	if err != nil {
		return err
	}

	daemon.Host = c.Dcompd.Host
	daemon.Port = c.Dcompd.Port

	auth := server.NewBasicAuth()
	daemon.SetAuth(auth)

	return nil
}
