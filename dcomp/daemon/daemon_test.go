package daemon

import (
	"bytes"
	"encoding/json"
	"io"
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

var submitTests = []request{
	{structs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", 201, "Create job"},
	{structs.JobDescription{"aaa", "nil", 1}, "jobs", "POST", 400, "create job - no server"},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", 400, "create job - nil struct"},
	{structs.JobDescription{}, "jobs", "POST", 400, "create job - empty struct"},
	{structs.JobDescription{}, "jobs/1", "POST", 404, "create job - wrong path"},
}

var getTests = []request{
	{structs.JobDescription{}, "jobs", "GET", 200, "Get all jobs"},
	{structs.JobDescription{}, "jobs/578359205e935a20adb39a18", "GET", 200, "Get existing job"},
	{structs.JobDescription{}, "jobs/1", "GET", 404, "Get non-existing job"},
	{structs.JobDescription{}, "job", "GET", 404, "get job - wrong path"},
}

func TestSubmitJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	for _, test := range submitTests {

		ts := server.CreateMockServer(&DBServer)
		if test.job.Script == "nil" {
			ts.Close()
		}
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.job.ImageName == "nil" {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
		ts.Close()
	}
}

func TestGetJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	for _, test := range getTests {

		ts := server.CreateMockServer(&DBServer)

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
		ts.Close()
	}
}
