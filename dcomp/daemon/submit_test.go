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

var submitRouteTests = []request{
	{structs.JobDescription{"ttt", "bbb", 1}, "jobs", "POST", http.StatusCreated, "Create job"},
	{structs.JobDescription{"hhh", "nil", 1}, "jobs", "POST", http.StatusBadRequest, "create job - no server"},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", http.StatusBadRequest, "create job - nil struct"},
	{structs.JobDescription{}, "jobs", "POST", http.StatusBadRequest, "create job - empty struct"},
	{structs.JobDescription{}, "jobs/1", "POST", http.StatusNotFound, "create job - wrong path"},
}

type submitRequest struct {
	job     structs.JobDescription
	answer  string
	message string
}

var submitTests = []submitRequest{
	{structs.JobDescription{"aaa", "bbb", 1}, "submittedimage", "Create job"},
	{structs.JobDescription{"aaa", "nil", 1}, "submittedimage", "no server"},
}

func TestRouteSubmitJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)

	for _, test := range submitRouteTests {

		ts := server.CreateMockServer(&dbServer)
		ts2 := server.CreateMockServer(&estimatorServer)
		defer ts2.Close()

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
		if w.Code == http.StatusCreated {
			assert.Contains(t, w.Body.String(), "submittedimage", test.message)
		}
		ts.Close()
	}
}

func TestSubmitJob(t *testing.T) {

	for _, test := range submitTests {

		ts := server.CreateMockServer(&dbServer)
		ts2 := server.CreateMockServer(&estimatorServer)
		defer ts2.Close()
		if test.job.Script == "nil" {
			ts.Close()
		}
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		job, err := submitJob(test.job)
		if test.job.Script == "nil" {
			assert.NotNil(t, err, "Should be error")
			continue
		}

		assert.Nil(t, err, "Should not be error")

		assert.Contains(t, job.ImageName, test.answer, test.message)
		ts.Close()
	}
}

func TestFindResource(t *testing.T) {
	for _, test := range submitTests {

		ts := server.CreateMockServer(&estimatorServer)
		if test.job.Script == "nil" {
			ts.Close()
		}

		prio, err := findResources(test.job)
		if test.job.Script == "nil" {
			assert.NotNil(t, err, "Should be error")
			continue
		}

		assert.Nil(t, err, "Should not be error")

		assert.Equal(t, 10, prio["Cloud"], test.message)
		ts.Close()
	}

}
