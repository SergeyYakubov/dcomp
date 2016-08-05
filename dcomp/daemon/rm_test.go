package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type request struct {
	job     structs.JobDescription
	path    string
	cmd     string
	answer  int
	message string
}

var getdeleteTests = []request{
	{structs.JobDescription{}, "jobs", "GET", 200, "Get all jobs"},
	{structs.JobDescription{}, "jobs/578359205e935a20adb39a18", "GET", 200, "Get existing job"},
	{structs.JobDescription{}, "jobs/1", "GET", 404, "Get non-existing job"},
	{structs.JobDescription{}, "job", "GET", 404, "get job - wrong path"},
	{structs.JobDescription{}, "jobs/578359205e935a20adb39a18", "DELETE", 200, "delete existing job"},
	{structs.JobDescription{}, "jobs/1", "DELETE", 404, "delete non-existing job"},
	{structs.JobDescription{}, "jobs", "DELETE", 404, "delete all jobs"},
	{structs.JobDescription{}, "job", "DELETE", 404, "delete job - wrong path"},
}

func TestRouteGetJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)

	for _, test := range getdeleteTests {

		ts := server.CreateMockServer(&dbServer)

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
		ts.Close()
	}
}