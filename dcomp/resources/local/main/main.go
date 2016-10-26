package main

import (
	"os"

	"log"

	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/resources/daemon"
	"stash.desy.de/scm/dc/main.git/dcomp/resources/local"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
	"stash.desy.de/scm/dc/main.git/dcomp/version"
)

type config struct {
	Daemon struct {
		Addr string
	}
	Database struct {
		Host string
		Port int
	}
}

func setConfiguration() (config, error) {
	fname := `/etc/dcomp/plugins/local.yaml`

	var c config

	err := utils.ReadYaml(fname, &c)
	return c, err
}

func main() {

	if ret := version.ShowVersion(os.Stdout, "dComp local resource plugin"); ret {
		return
	}

	c, err := setConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	var dbsrv server.Server
	dbsrv.Host = c.Database.Host
	dbsrv.Port = c.Database.Port
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugin")
	addr := c.Daemon.Addr

	var res = new(local.Resource)

	daemon.Start(res, db, addr)
	return
}
