package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type request struct {
	job    structs.JobDescription
	path   string
	cmd    string
	msg    string
	answer int
}

var submitTests = []request{
	{structs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", "post normal job", 201},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", "post with nil body", 400},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs", "GET", "get all jobs, nil body", 200},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs/1", "GET", "get job 1", 200},
	{structs.JobDescription{}, "jobs", "POST", "post empty structure", 400},
	{structs.JobDescription{}, "jobs", "GET", "get all jobs, empty structure", 200},
	{structs.JobDescription{}, "jobs/1", "GET", "get job 1", 200},
	{structs.JobDescription{}, "jobs/1", "POST", "post job 1", 404},
	{structs.JobDescription{}, "job", "GET", "wrong path", 404},
}

func TestSubmitJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	if err := database.CreateMock(); err != nil {
		t.Error("Create database" + err.Error())
		return
	}

	if err := database.SetServerConfiguration(); err != nil {
		t.Error("Set server config" + err.Error())
		return

	}

	if err := database.Connect(); err != nil {
		t.Error("Connect to the database " + err.Error())
		return
	}

	for _, test := range submitTests {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.job.ImageName == "nil" {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8001/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.msg)
	}
}
