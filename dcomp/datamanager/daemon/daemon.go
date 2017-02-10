// Package for estimator daemon
package daemon

import (
	"log"
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/datamanager/internal/cachedatabase"
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
	Database struct {
		Host string
		Port int
	}
}

var settings config

var db cachedatabase.Agent

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration(configFile string) error {

	return utils.ReadYaml(configFile, &settings)

}

func Start(configFile string) {

	mux := utils.NewRouter(listRoutes)
	if err := setDaemonConfiguration(configFile); err != nil {
		log.Fatal(err)
	}

	db = new(cachedatabase.SqlDatabase)
	srv := server.Server{Host: settings.Database.Host, Port: settings.Database.Port}
	db.SetServer(&srv)

	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Fatal(http.ListenAndServeTLS(settings.Daemon.Addr,
		settings.Daemon.Certfile, settings.Daemon.Keyfile,
		server.ProcessJWTAuth(mux.ServeHTTP, settings.Daemon.Key)))
}
