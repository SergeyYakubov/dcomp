package daemon

import (
	"fmt"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var estimatorServer server.Server
var resources map[string]structs.Resource
var dbServer server.Server
var addr string

type config struct {
	Daemon struct {
		Addr string
	}
	Database struct {
		Host string
		Port int
	}
	Estimator struct {
		Host string
		Port int
	}
	Plugins []struct {
		Name string
		Host string
		Port int
	}
}

func setConfiguration() error {

	fname := `/etc/dcomp/dcompd.yaml`

	var c config
	if err := utils.ReadYaml(fname, &c); err != nil {
		fmt.Println(err)
		return err
	}

	dbServer.Host = c.Database.Host
	dbServer.Port = c.Database.Port

	estimatorServer.Host = c.Estimator.Host
	estimatorServer.Port = c.Estimator.Port

	addr = c.Daemon.Addr

	resources = make(map[string]structs.Resource)

	// add plugins found in config file
	for _, p := range c.Plugins {
		resources[p.Name] = structs.Resource{Server: server.Server{p.Host, p.Port}}
	}

	return nil
}

func connectDb(name string) (db database.Agent, err error) {

	db = new(database.Mongodb)
	if err := setConfiguration(); err != nil {
		return nil, err
	}

	db.SetServer(&dbServer)
	db.SetDefaults("daemondbd")

	if err = db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}
