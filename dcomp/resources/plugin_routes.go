package resources

import "stash.desy.de/scm/dc/main.git/dcomp/utils"

func (p *plugin) setRoutes() {
	p.ListRoutes = utils.Routes{
		utils.Route{
			"SubmitJob",
			"POST",
			"/jobs/",
			p.SubmitJob,
		},
	}
}
