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
var authServer server.Server
var settings struct {
	Addr     string
	Certfile string
	Keyfile  string
}

type HostInfo struct {
	Host string
	Port int
	Key  string
}

type config struct {
	Daemon struct {
		Addr     string
		Certfile string
		Keyfile  string
	}
	Database      HostInfo
	Estimator     HostInfo
	Authorization HostInfo

	Plugins []struct {
		Host        string
		Port        int
		Key         string
		Name        string
		DataManager struct {
			Host string
			Port int
			Key  string
		}
	}
}

func setHostInfo(srv *server.Server, h HostInfo) {
	srv.Host = h.Host
	srv.Port = h.Port
	auth := server.NewHMACAuth(h.Key)
	srv.SetAuth(auth)

}

var configFile = `/etc/dcomp/conf/dcompd.yaml`

func setConfiguration() error {

	var c config
	if err := utils.ReadYaml(configFile, &c); err != nil {
		fmt.Println(err)
		return err
	}

	setHostInfo(&dbServer, c.Database)
	setHostInfo(&estimatorServer, c.Estimator)
	setHostInfo(&authServer, c.Authorization)
	authServer.Tls = true

	settings.Addr = c.Daemon.Addr
	settings.Certfile = c.Daemon.Certfile
	settings.Keyfile = c.Daemon.Keyfile

	resources = make(map[string]structs.Resource)

	// add plugins found in config file
	for _, p := range c.Plugins {
		s := server.Server{Host: p.Host, Port: p.Port}
		auth := server.NewHMACAuth(p.Key)
		s.SetAuth(auth)
		dm := server.Server{Host: p.DataManager.Host, Port: p.DataManager.Port}
		auth2 := server.NewJWTAuth(p.DataManager.Key)
		dm.SetAuth(auth2)
		dm.Tls = true
		resources[p.Name] = structs.Resource{Server: s, DataManager: dm}
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
