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

var submitRouteTests = []request{
	{structs.JobDescription{},
		"jobs/578359205e935a20adb39a18", "POST", http.StatusCreated, "Create job"},
}

func TestRouteSubmitReleaseJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitRouteTests {

		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		r := resources["local"]
		r.Server = srv
		resources["local"] = r
		resources["mock"] = r

		defer ts2.Close()

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.answer, w.Code, test.message)
		if w.Code == http.StatusCreated {
			assert.Contains(t, w.Body.String(), "578359205e935a20adb39a18", test.message)
		}
		if w.Code == http.StatusAccepted {
			assert.Contains(t, w.Body.String(), "Bearer", test.message)
		}

		ts3.Close()
	}
}

type submitRequest struct {
	job     structs.JobDescription
	answer  string
	status  int
	message string
}

var submitTests = []submitRequest{
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: "local"}, "578359205e935a20adb39a18",
		structs.StatusSubmitted, "Create job"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: ""}, "578359205e935a20adb39a18",
		structs.StatusSubmitted, "Create job"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: "local",
		FilesToUpload: structs.TransferFiles{
			{"jhjh", "assd"},
			{"jhjh", "assd"},
		}}, "578359205e935a20adb39a18", structs.StatusWaitData, "Wait for files"},
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: 1, Resource: "local"}, "available",
		structs.StatusError, "Create job"},
}

func TestSubmitJob(t *testing.T) {

	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitTests {
		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		resources["local"] = structs.Resource{Server: srv}

		defer ts2.Close()

		if test.job.ImageName == "nil" {
			ts3.Close()
		}

		job, err := trySubmitJob("", test.job)
		if test.job.Script == "nil" || test.job.ImageName == "nil" {
			assert.NotNil(t, err, "Should be error")
			if err != nil {
				assert.Contains(t, err.Error(), test.answer, test.message)
			}
			continue
		}

		assert.Nil(t, err, "Should not be error")

		assert.Equal(t, test.status, job.Status, test.message)
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
		if test.job.Resource != "" {
			assert.Equal(t, 100, prio[test.job.Resource], test.message)
		} else {
			assert.Equal(t, 10, prio["cloud"], test.message)
		}
		ts.Close()
	}
}
