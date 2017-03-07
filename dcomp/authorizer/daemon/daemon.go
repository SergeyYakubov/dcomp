// Package for authorizer daemon
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
	Tokenduration int
	Ldap          struct {
		Host   string
		BaseDn string
	}
	Authorization []string
}

func (c *config) authorizationAllowed(atype string) bool {
	return utils.StringInArray(atype, c.Authorization)
}

var configFile = `/etc/dcomp/conf/dcompauthd.yaml`
var c config

// setDaemonConfiguration reads configuration file with daemon location
func setDaemonConfiguration() error {

	err := utils.ReadYaml(configFile, &c)
	return err

}

var gssAPIContext *utils.Context

func Start() {

	mux := utils.NewRouter(listRoutes)
	var err error
	err = setDaemonConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	if utils.StringInArray("Negotiate", c.Authorization) {
		gssAPIContext, err = utils.PrepareSeverGSSAPIContext("dcomp")

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(http.ListenAndServeTLS(c.Daemon.Addr, c.Daemon.Certfile, c.Daemon.Keyfile,
		server.ProcessHMACAuth(mux.ServeHTTP, c.Daemon.Key)))

}
