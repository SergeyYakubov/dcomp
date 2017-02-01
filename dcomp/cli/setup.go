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

var configFile = `/etc/dcomp/conf/dcomp.yaml`

// SetDaemonConfiguration reads configuration file with daemon location
func SetDaemonConfiguration() error {

	var c config

	err := utils.ReadYaml(configFile, &c)

	if err != nil {
		return err
	}

	daemon.Host = c.Dcompd.Host
	daemon.Port = c.Dcompd.Port

	//	auth := server.NewBasicAuth()
	auth := server.NewGSSAPIAuth()
	daemon.SetAuth(auth)
	daemon.Tls = true

	return nil
}
