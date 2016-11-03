package main

import (
	"os"

	"log"

	"github.com/pkg/errors"
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
	BaseDir string
}

func (c *config) check() error {

	src, err := os.Stat(c.BaseDir)
	if err != nil {
		return err
	}

	if !src.IsDir() {
		err := errors.New(c.BaseDir + " is not a directory")
		return err
	}
	return nil
}

var configFileName = `/etc/dcomp/plugins/local.yaml`

func setConfiguration() (config, error) {

	var c config

	if err := utils.ReadYaml(configFileName, &c); err != nil {
		return c, err
	}

	err := c.check()
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

	res.Basedir = c.BaseDir
	daemon.Start(res, db, addr)
	return
}
