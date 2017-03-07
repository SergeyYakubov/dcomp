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
		"Login",
		"GET",
		"/login/",
		routeLogin,
	},

	utils.Route{
		"PatchJob",
		"PATCH",
		"/jobs/{jobID}/",
		routePatchJob,
	},
	utils.Route{
		"GetJob",
		"POST",
		"/jobs/{jobID}/",
		routeReleaseJob,
	},
	utils.Route{
		"GetJobFiles",
		"GET",
		"/jobfile/{jobID}/",
		SendJWTToken,
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
