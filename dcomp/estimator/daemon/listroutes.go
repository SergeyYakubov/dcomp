package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var ListRoutes = utils.Routes{
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		EstimateJob,
	},
}
