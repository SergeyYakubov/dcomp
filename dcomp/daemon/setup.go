package daemon

import (
	"fmt"
	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
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
		Key  string
	}
	Estimator struct {
		Host string
		Port int
		Key  string
	}
	Plugins []struct {
		Name string
		Host string
		Port int
		Key  string
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
	dbServer.Key = c.Database.Key

	estimatorServer.Host = c.Estimator.Host
	estimatorServer.Port = c.Estimator.Port
	estimatorServer.Key = c.Estimator.Key

	addr = c.Daemon.Addr

	resources = make(map[string]structs.Resource)

	// add plugins found in config file
	for _, p := range c.Plugins {
		resources[p.Name] = structs.Resource{Server: server.Server{p.Host, p.Port, p.Key}}
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
