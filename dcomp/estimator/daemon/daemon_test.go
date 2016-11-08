package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/dcomp/dcomp/structs"
	"github.com/dcomp/dcomp/utils"
)

type request struct {
	job     structs.JobDescription
	path    string
	cmd     string
	answer  int
	message string
	prio    structs.ResourcePrio
}

var estimateTests = []request{
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Local: true}, "estimations", "POST", http.StatusOK, "estimations job batch",
		structs.ResourcePrio{"Local": 100, "HPC": 1, "Batch": 10, "Cloud": 0}},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 8}, "estimations", "POST", http.StatusOK, "estimations job batch/hpc",
		structs.ResourcePrio{"HPC": 5, "Batch": 5, "Cloud": 0}},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 80}, "estimations", "POST", http.StatusOK, "estimations job hpc",
		structs.ResourcePrio{"HPC": 10, "Batch": 0, "Cloud": 0}},
	{structs.JobDescription{}, "estimations", "POST", http.StatusBadRequest, "estimations job - empty struct",
		structs.ResourcePrio{}},
	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: -1}, "estimations", "POST", http.StatusBadRequest, "create job - nil struct",
		structs.ResourcePrio{}},
	{structs.JobDescription{}, "estimations/1", "POST", http.StatusNotFound, "estimations job - wrong path",
		structs.ResourcePrio{}},
}

func TestEstimateJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)

	for _, test := range estimateTests {

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.job); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.job.ImageName == "nil" {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8002/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
		if w.Code == http.StatusOK {
			var prio structs.ResourcePrio
			if err := json.NewDecoder(w.Body).Decode(&prio); err != nil {
				t.Fail()
			}
			assert.Equal(t, test.prio, prio, test.message)

		}
	}
}
