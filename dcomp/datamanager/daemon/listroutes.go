package daemon

import (
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"ReceiveFile",
		"POST",
		"/jobfile/{jobID}/",
		routeReceiveJobFile,
	},
	utils.Route{
		"SendFile",
		"GET",
		"/jobfile/{jobID}/",
		routeSendJobFile,
	},
}
