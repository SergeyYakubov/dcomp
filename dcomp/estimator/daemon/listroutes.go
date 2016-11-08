package daemon

import (
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		routeEstimateJob,
	},
}
