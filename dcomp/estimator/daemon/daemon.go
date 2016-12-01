// Package for estimator daemon
package daemon

import (
	"log"
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type config struct {
	Daemon struct {
		Addr string
		Key  string
	}
}

var configFile = `/etc/dcomp/conf/dcompestd.yaml`

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration() (config, error) {

	var c config

	err := utils.ReadYaml(configFile, &c)
	return c, err

}

var c config

func Start() {

	mux := utils.NewRouter(listRoutes)
	var err error
	c, err = setDaemonConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(c.Daemon.Addr, server.ProcessHMACAuth(mux.ServeHTTP, c.Daemon.Key)))
}
