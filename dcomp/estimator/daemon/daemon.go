// Package for estimator daemon
package daemon

import (
	"log"
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

type config struct {
	Daemon struct {
		Addr string
		Key  string
	}
}

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration() (config, error) {

	fname := `/etc/dcomp/dcompestd.yaml`

	var c config

	err := utils.ReadYaml(fname, &c)
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
