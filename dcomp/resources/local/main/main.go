package main

import (
	"os"

	"log"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/resources/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/resources/local"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

type config struct {
	Daemon struct {
		Addr string
		Key  string
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
	key := c.Daemon.Key

	var res = new(local.Resource)

	res.Basedir = c.BaseDir
	daemon.Start(res, db, addr, key)
	return
}
