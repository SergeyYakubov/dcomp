package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/common_structs"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var submitTests = []struct {
	job        commonStructs.JobDescription
	path       string
	cmd        string
	serverresp string
	answer     int
}{
	{commonStructs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", "ok", 201},
	{commonStructs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", "badreq", 400},
	{commonStructs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", "empty", 201},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", "ok", 400},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs", "GET", "ok", 200},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs/1", "GET", "ok", 200},
	{commonStructs.JobDescription{}, "jobs", "POST", "ok", 400},
	{commonStructs.JobDescription{}, "jobs", "GET", "ok", 200},
	{commonStructs.JobDescription{}, "jobs/1", "GET", "ok", 200},
	{commonStructs.JobDescription{}, "jobs/1", "POST", "ok", 404},
	{commonStructs.JobDescription{}, "job", "GET", "ok", 404},
}

func TestSubmitJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	for _, test := range submitTests {
		ts := server.CreateMockServer(&DBServer, test.serverresp)
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
		assert.Equal(t, test.answer, w.Code)
		ts.Close()
	}
}
