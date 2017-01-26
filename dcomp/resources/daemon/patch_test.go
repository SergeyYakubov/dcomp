package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"

	"github.com/sergeyyakubov/dcomp/dcomp/resources/mock"
	"github.com/sergeyyakubov/dcomp/dcomp/server"

	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type patchRequest struct {
	patch      structs.PatchJob
	path       string
	cmd        string
	answercode int
	message    string
}

var patchTests = []patchRequest{
	{structs.PatchJob{Status: structs.StatusFinished}, "jobs/578359205e935a20adb39a18", "PATCH",
		http.StatusOK, "kill job"},
	{structs.PatchJob{Status: structs.StatusFinished}, "jobs/578359205e935a20adb39a19", "PATCH",
		http.StatusBadRequest, "job not found"},
	{structs.PatchJob{Status: structs.StatusError}, "jobs/578359205e935a20adb39a18", "PATCH",
		http.StatusBadRequest, "empty status"},
	{structs.PatchJob{Status: structs.StatusFailed}, "jobs/578359205e935a20adb39a18", "PATCH",
		http.StatusBadRequest, "wrong status"},
	{structs.PatchJob{Status: structs.StatusFinished}, "jobs", "PATCH",
		http.StatusNotFound, "no job id "},
}

func TestPatchJob(t *testing.T) {
	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mockdatabase)
	db.SetServer(&dbsrv)
	db.Connect()
	resource = new(mock.MockResource)
	resource.SetDb(db)

	mux := utils.NewRouter(listRoutes)
	for _, test := range patchTests {

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.patch); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.patch.Status == structs.StatusError {
			reader = nil
		}

		req, err := http.NewRequest(test.cmd, "http://localhost:8002/"+test.path+"/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answercode, w.Code, test.message)
	}

}
