package daemon

import (
	"github.com/dcomp/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		routeEstimateJob,
	},
}
