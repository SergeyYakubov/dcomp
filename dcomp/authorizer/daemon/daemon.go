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
		Addr     string
		Key      string
		Certfile string
		Keyfile  string
	}
}

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration() (config, error) {

	fname := `/etc/dcomp/conf/dcompauthd.yaml`

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
	log.Fatal(http.ListenAndServeTLS(c.Daemon.Addr, c.Daemon.Certfile, c.Daemon.Keyfile,
		server.ProcessHMACAuth(mux.ServeHTTP, c.Daemon.Key)))

}
