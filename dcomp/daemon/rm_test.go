package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

type request struct {
	job     structs.JobDescription
	path    string
	cmd     string
	answer  int
	message string
}

var getdeleteTests = []request{
	//	{structs.JobDescription{}, "jobs", "GET", 200, "Get all jobs"},
	//	{structs.JobDescription{}, "jobs/578359205e935a20adb39a18", "GET", 200, "Get existing job"},
	//	{structs.JobDescription{}, "jobs/1", "GET", 404, "Get non-existing job"},
	//	{structs.JobDescription{}, "job", "GET", 404, "get job - wrong path"},
	{structs.JobDescription{}, "jobs/578359205e935a20adb39a18", "DELETE", http.StatusNoContent, "delete existing job"},
	{structs.JobDescription{}, "jobs/578359205e935a20adb39a19", "DELETE", 404, "delete non-existing job"},
	{structs.JobDescription{}, "jobs", "DELETE", 404, "delete all jobs"},
	{structs.JobDescription{}, "job", "DELETE", 404, "delete job - wrong path"},
}

func TestRouteDeleteJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	var srv server.Server
	ts3 := server.CreateMockServer(&srv)
	defer ts3.Close()
	resources = make(map[string]structs.Resource)
	resources["mock"] = structs.Resource{Server: srv}

	for _, test := range getdeleteTests {

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
	}
}
