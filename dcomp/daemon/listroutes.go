package daemon

import (
	"github.com/dcomp/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		routeGetAllJobs,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}/",
		routeGetJob,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		routeSubmitJob,
	},
	utils.Route{
		"DeleteJob",
		"DELETE",
		"/jobs/{jobID}/",
		routeDeleteJob,
	},
}
