package resources

import (
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type plugin struct {
	database   database.Agent
	resource   Resource
	ListRoutes utils.Routes
}

func NewPlugin(r Resource, db database.Agent) *plugin {
	p := new(plugin)
	r.SetUpdateStatusCmd(p.UpdateJobStatus)
	p.resource = r
	p.database = db
	p.setRoutes()
	return p
}
