package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"net/http/httptest"

	"stash.desy.de/scm/dc/main.git/dcomp/resources/mock"

	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

type jobDeleteRequest struct {
	path       string
	cmd        string
	answercode int
	answer     string
	message    string
}

var jobDeleteTests = []jobDeleteRequest{
	{"jobs/578359205e935a20adb39a18", "DELETE", http.StatusOK, "200", "delete job ok "},
	{"jobs/578359205e935a20adb39a19", "DELETE", http.StatusBadRequest, "12345", "no job found "},
	{"jobs/1", "DELETE", http.StatusBadRequest, "12345", "no job found "},
	{"jobs", "DELETE", http.StatusNotFound, "12345", "no job found "},
	{"aaa", "DELETE", http.StatusNotFound, "12345", "no job found "},
}

func TestDeleteJob(t *testing.T) {
	resource = new(mock.MockResource)

	mux := utils.NewRouter(listRoutes)
	for _, test := range jobDeleteTests {

		req, err := http.NewRequest(test.cmd, "http://localhost:8002/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answercode, w.Code, test.message)
	}

}
