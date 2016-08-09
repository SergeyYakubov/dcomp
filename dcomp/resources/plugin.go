package resources

import (
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type plugin struct {
	database   database.Agent
	resource   Resource
	listRoutes utils.Routes
}

func NewPlugin(r Resource) *plugin {
	p := new(plugin)
	p.resource = r
	p.setRoutes()
	return p
}

func (p *plugin) CloseDatabase() {
	if p.database != nil {
		p.database.Close()
	}
}

func (p *plugin) InitializeDatabase(db database.Agent, srv server.Server, defaults ...interface{}) error {
	p.database = db
	p.database.SetServer(&srv)
	p.database.SetDefaults(defaults)
	return p.database.Connect()
}
