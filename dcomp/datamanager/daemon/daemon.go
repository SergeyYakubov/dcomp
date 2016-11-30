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
	Dcompd struct {
		Host string
		Port int
	}
	Resource struct {
		BaseDir string
	}
}

var settings config

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration(configFile string) error {

	return utils.ReadYaml(configFile, &settings)

}

func Start(configFile string) {

	mux := utils.NewRouter(listRoutes)
	if err := setDaemonConfiguration(configFile); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServeTLS(settings.Daemon.Addr,
		settings.Daemon.Certfile, settings.Daemon.Keyfile,
		server.ProcessJWTAuth(mux.ServeHTTP, settings.Daemon.Key)))
}
