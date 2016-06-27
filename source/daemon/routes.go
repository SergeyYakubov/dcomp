package daemon

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Routes []Route

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range ListRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var ListRoutes = Routes{
	Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		GetAllJobs,
	},
	Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}",
		GetJob,
	},
	Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		SubmitJob,
	},
}
