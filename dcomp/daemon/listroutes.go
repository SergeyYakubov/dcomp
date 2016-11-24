package daemon

import (
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
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
		"GetJob",
		"POST",
		"/jobs/{jobID}/",
		routeReleaseJob,
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
