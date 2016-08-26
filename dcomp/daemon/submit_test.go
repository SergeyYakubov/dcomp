package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var submitRouteTests = []request{
	{structs.JobDescription{ImageName: "ttt", Script: "bbb", NCPUs: 1, Local: true}, "jobs", "POST", http.StatusCreated, "Create job"},
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: -1}, "jobs", "POST", http.StatusBadRequest, "create job - nil struct"},
	{structs.JobDescription{}, "jobs", "POST", http.StatusBadRequest, "create job - empty struct"},
	{structs.JobDescription{}, "jobs/1", "POST", http.StatusNotFound, "create job - wrong path"},
}

type submitRequest struct {
	job     structs.JobDescription
	answer  string
	message string
}

var submitTests = []submitRequest{
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Local: true}, "578359205e935a20adb39a18", "Create job"},
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: 1, Local: true}, "available", "Create job"},
}

func TestRouteSubmitJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	initialize()
	db = new(database.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitRouteTests {

		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		resources["Local"] = structs.Resource{Server: srv}

		defer ts2.Close()

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
			assert.Contains(t, w.Body.String(), "578359205e935a20adb39a18", test.message)
		}
		ts3.Close()
	}
}

func TestSubmitJob(t *testing.T) {

	initialize()
	db = new(database.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitTests {
		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		resources["Local"] = structs.Resource{Server: srv}

		defer ts2.Close()

		if test.job.ImageName == "nil" {
			ts3.Close()
		}

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		job, err := submitJob(test.job)
		if test.job.Script == "nil" || test.job.ImageName == "nil" {
			assert.NotNil(t, err, "Should be error")
			assert.Contains(t, err.Error(), test.answer, test.message)
			continue
		}

		assert.Nil(t, err, "Should not be error")

		assert.Contains(t, job.Id, test.answer, test.message)
		ts3.Close()
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
