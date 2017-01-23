package main

import (
	"os"

	"log"

	"flag"
	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/resources/cluster"
	"github.com/sergeyyakubov/dcomp/dcomp/resources/daemon"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/sergeyyakubov/dcomp/dcomp/version"
)

type config struct {
	PluginName string
	Daemon     struct {
		Addr string
		Key  string
	}
	Database struct {
		Host string
		Port int
	}
	BaseDir     string
	TemplateDir string
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

func setConfiguration(configFileName string) (config, error) {

	var c config

	if err := utils.ReadYaml(configFileName, &c); err != nil {
		return c, err
	}

	err := c.check()
	return c, err

}

func main() {

	if ret := version.ShowVersion(os.Stdout, "dComp cluster plugin"); ret {
		return
	}

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("configuration file not set. Usage: dcompclusterpd <config file name>")
	}

	configFileName := flag.Arg(0)

	c, err := setConfiguration(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	var dbsrv server.Server
	dbsrv.Host = c.Database.Host
	dbsrv.Port = c.Database.Port
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults(c.PluginName + "plugin")
	addr := c.Daemon.Addr
	key := c.Daemon.Key

	var res = new(cluster.Resource)

	res.Basedir = c.BaseDir
	res.TemplateDir = c.TemplateDir
	res.Name = c.PluginName + " plugin"

	daemon.Start(res, db, addr, key)
	return
}
