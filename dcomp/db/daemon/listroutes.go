package daemon

import (
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		getAllJobs,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}/",
		getJob,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		submitJob,
	},
	utils.Route{
		"DeleteJob",
		"DELETE",
		"/jobs/{jobID}/",
		deleteJob,
	},
}
