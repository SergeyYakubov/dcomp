package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		routeEstimateJob,
	},
}
