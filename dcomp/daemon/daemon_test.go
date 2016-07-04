package daemon

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"stash.desy.de/scm/dc/main.git/dcomp/common_structs"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
	"testing"
)

var submitTests = []struct {
	job    commonStructs.JobDescription
	path   string
	cmd    string
	answer int
}{
	{commonStructs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", 200},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", 400},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs", "GET", 200},
	{commonStructs.JobDescription{"nil", "bbb", -1}, "jobs/1", "GET", 200},
	{commonStructs.JobDescription{}, "jobs", "POST", 400},
	{commonStructs.JobDescription{}, "jobs", "GET", 200},
	{commonStructs.JobDescription{}, "jobs/1", "GET", 200},
	{commonStructs.JobDescription{}, "jobs/1", "POST", 404},
	{commonStructs.JobDescription{}, "job", "GET", 404},
}

func TestSubmitJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	ts := server.CreateMockServer(&DBServer)
	defer ts.Close()

	for _, test := range submitTests {

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
	}

}
