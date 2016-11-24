package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

var submitRouteTests = []request{
	/*{structs.JobDescription{ImageName: "ttt", Script: "bbb", NCPUs: 1, Local: true}, "jobs", "POST", http.StatusCreated, "Create job"},
	{structs.JobDescription{ImageName: "ttt", Script: "bbb", NCPUs: 1, Local: true,
		FilesToUpload: structs.TransferFiles{
			{"jhjh", "assd"},
			{"jhjh", "assd"}},
	}, "jobs", "POST", http.StatusAccepted, "Create job,wait files"},
	{structs.JobDescription{ImageName: "noauth", Script: "bbb", NCPUs: 1, Local: true,
		FilesToUpload: structs.TransferFiles{
			{"jhjh", "assd"},
			{"jhjh", "assd"}},
	}, "jobs", "POST", http.StatusInternalServerError, "Create job,wait files, no auth context"},

	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: -1}, "jobs", "POST", http.StatusBadRequest, "create job - nil struct"},
	{structs.JobDescription{}, "jobs", "POST", http.StatusBadRequest, "create job - empty struct"},
	{structs.JobDescription{}, "jobs/1", "POST", http.StatusBadRequest, "create job - wrong path"},*/
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: 1, Local: true},
		"jobs/578359205e935a20adb39a18", "POST", http.StatusCreated, "Create job"},
}

func TestRouteSubmitReleaseJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	setConfiguration()
	db = new(database.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitRouteTests {

		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		r := resources["Local"]
		r.Server = srv
		resources["Local"] = r
		resources["mock"] = r

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

		if test.job.ImageName != "noauth" {
			resp := server.AuthorizationResponce{"testuser"}
			ctx := context.WithValue(req.Context(), "authorizationResponce", &resp)
			req = req.WithContext(ctx)
		}

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
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Local: true}, "578359205e935a20adb39a18",
		structs.StatusSubmitted, "Create job"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Local: true,
		FilesToUpload: structs.TransferFiles{
			{"jhjh", "assd"},
			{"jhjh", "assd"},
		}}, "578359205e935a20adb39a18", structs.StatusWaitData, "Wait for files"},
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: 1, Local: true}, "available",
		structs.StatusError, "Create job"},
}

func TestSubmitJob(t *testing.T) {

	setConfiguration()
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

		job, err := trySubmitJob(test.job)
		if test.job.Script == "nil" || test.job.ImageName == "nil" {
			assert.NotNil(t, err, "Should be error")
			assert.Contains(t, err.Error(), test.answer, test.message)
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

		assert.Equal(t, 10, prio["Cloud"], test.message)
		ts.Close()
	}
}
