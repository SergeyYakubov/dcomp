// +build !release

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"stash.desy.de/scm/dc/utils"
)

func MockFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"ImageName":"ddd","Script":"aaa","NCPUs":1,"Id":1,"Status":1}`)
}

var listRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		nil,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}",
		nil,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		nil,
	},
}

func CreateMockServer(srv *Srv) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(MockFunc))
	srv.ParseUrl(ts.URL)
	return ts
}
