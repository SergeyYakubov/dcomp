package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"

	"stash.desy.de/scm/dc/main.git/dcomp/resources/mock"
	"stash.desy.de/scm/dc/main.git/dcomp/server"

	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type request struct {
	job        structs.JobInfo
	path       string
	cmd        string
	answercode int
	answer     string
	message    string
}

var submitTests = []request{
	{structs.JobInfo{JobDescription: structs.JobDescription{"image", "script", 20},
		Id: "578359205e935a20adb39a18"}, "jobs", "POST", http.StatusCreated, "12345", "submit job"},
	{structs.JobInfo{JobDescription: structs.JobDescription{"image", "script", 20},
		Id: "578359205e935a20adb39a18"}, "job", "POST", http.StatusNotFound, "12345", "wrong path"},
	{structs.JobInfo{JobDescription: structs.JobDescription{"nil", "script", 20},
		Id: "578359205e935a20adb39a18"}, "jobs", "POST", http.StatusBadRequest, "12345", "wrong input"},
	{structs.JobInfo{JobDescription: structs.JobDescription{"errorsubmit", "script", 20},
		Id: "578359205e935a20adb39a18"}, "jobs", "POST", http.StatusInternalServerError, "error", "error from resource"},
	{structs.JobInfo{JobDescription: structs.JobDescription{"image", "script", 20},
		Id: "give error"}, "jobs", "POST", http.StatusInternalServerError, "error create", "error from database"},
}

func TestSubmitJob(t *testing.T) {
	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mockdatabase)
	db.SetServer(&dbsrv)
	db.Connect()
	p := NewPlugin(new(mock.MockResource), db)

	mux := utils.NewRouter(p.ListRoutes)
	for _, test := range submitTests {

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
		assert.Equal(t, test.answercode, w.Code, test.message)
		if w.Code == http.StatusOK {
			assert.Contains(t, w.Body, test.answer, test.message)
		}
	}

}
