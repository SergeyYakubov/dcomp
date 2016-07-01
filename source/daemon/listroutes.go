package daemon

import (
	"stash.desy.de/scm/dc/utils"
)

var ListRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		GetAllJobs,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}/",
		GetJob,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		SubmitJob,
	},
}
