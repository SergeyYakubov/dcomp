package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type submitrequest struct {
	job    structs.JobDescription
	path   string
	cmd    string
	msg    string
	answer int
}

type getrequest struct {
	path   string
	cmd    string
	body   string
	msg    string
	answer int
}

var submitTests = []submitrequest{
	{structs.JobDescription{"aaa", "bbb", 1}, "jobs", "POST", "post normal job", 201},
	{structs.JobDescription{"nil", "bbb", -1}, "jobs", "POST", "post with nil body", 400},
	{structs.JobDescription{}, "jobs", "POST", "post empty structure", 400},
	{structs.JobDescription{}, "jobs/1", "POST", "post job 1", 404},
}

var getTests = []getrequest{
	//{"jobs", "GET", "get all jobs", 200},
	{"jobs/578359205e935a20adb39a18", "GET", "578359205e935a20adb39a18", "get job 1", 200},
	{"jobs/578359205e935a20adb39a19", "GET", "not found", "job not exist", 404},
	{"jobs/2", "GET", "cannot", "wrong format", 400},
	{"job", "GET", "not found", "wrong path", 404},
	{"jobs/578359205e935a20adb39a18", "DELETE", "", "delete job", 200},
	{"jobs/578359205e935a20adb39a19", "DELETE", "not found", "job not exist", 400},
	{"jobs", "DELETE", "not found", "delete all jobs", 404},
	{"jobs/2", "DELETE", "cannot", "wrong format", 400},
	{"job", "DELETE", "not found", "wrong path", 404},
}

func preparedatabase() error {
	if err := database.CreateMock(); err != nil {
		return errors.New("Create database" + err.Error())
	}

	if err := database.SetServerConfiguration(); err != nil {
		return errors.New("Set server config" + err.Error())
	}

	if err := database.Connect(); err != nil {
		return errors.New("Connect to the database " + err.Error())
	}
	return nil
}

func TestSubmitJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	if err := preparedatabase(); err != nil {
		t.Error("Create database" + err.Error())
		return
	}
	defer database.Close()

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

func TestGetDeleteJob(t *testing.T) {
	mux := utils.NewRouter(ListRoutes)

	if err := preparedatabase(); err != nil {
		t.Error("Create database" + err.Error())
		return
	}
	defer database.Close()

	s := structs.JobInfo{structs.JobDescription{}, "dummyid", 1}
	database.CreateRecord(&s)

	for _, test := range getTests {

		req, err := http.NewRequest(test.cmd, "http://localhost:8001/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answer, w.Code, test.msg)
		assert.Contains(t, w.Body.String(), test.body, test.msg)
	}

}
