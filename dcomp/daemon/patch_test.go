package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
	"io"
)

type patchRequest struct {
	patch   structs.PatchJob
	path    string
	cmd     string
	answer  int
	message string
}

var patchTests = []patchRequest{
	{structs.PatchJob{Status: structs.StatusFinished}, "jobs/578359205e935a20adb39a18", "PATCH", http.StatusNoContent, "patch existing job"},
	{structs.PatchJob{Status: structs.StatusFinished}, "jobs/578359205e935a20adb39a19", "PATCH", http.StatusNotFound, "patch non-existing job"},
	{structs.PatchJob{Status: structs.StatusError}, "jobs/578359205e935a20adb39a18", "PATCH", http.StatusBadRequest, "patch no data given"},
	{structs.PatchJob{}, "jobs", "PATCH", http.StatusNotFound, "jobs id not given"},
	{structs.PatchJob{}, "job", "PATCH", http.StatusNotFound, "patch job - wrong path"},
}

func TestRoutePatchJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	var srv server.Server
	ts3 := server.CreateMockServer(&srv)
	defer ts3.Close()
	resources = make(map[string]structs.Resource)
	resources["mock"] = structs.Resource{Server: srv}

	for _, test := range patchTests {

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.patch); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.patch.Status == structs.StatusError {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.message)
	}
}
