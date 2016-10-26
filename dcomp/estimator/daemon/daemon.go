// Package for estimator daemon
package daemon

import (
	"log"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type config struct {
	Daemon struct {
		Addr string
	}
}

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration() (config, error) {

	fname := `/etc/dcomp/dcompestd.yaml`

	var c config

	err := utils.ReadYaml(fname, &c)
	return c, err

}

func Start() {

	mux := utils.NewRouter(listRoutes)
	c, err := setDaemonConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(c.Daemon.Addr, mux))
}
