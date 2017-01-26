package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"net/http/httptest"

	"github.com/sergeyyakubov/dcomp/dcomp/resources/mock"

	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type jobInfoRequest struct {
	path       string
	querry     string
	cmd        string
	answercode int
	answer     string
	message    string
}

var JobInfoTests = []jobInfoRequest{
	{"jobs/578359205e935a20adb39a18", "", "GET", http.StatusOK, "102", "get job info "},
	{"jobs/578359205e935a20adb39a18", "?log=true", "GET", http.StatusOK, "hello", "get job info "},
	{"jobs/578359205e935a20adb39a18", "?log=true&compress=true", "GET", http.StatusOK,
		utils.CompressString("hello"), "get job info "},
	{"jobs/578359205e935a20adb39a19", "", "GET", http.StatusNotFound, "12345", "no job found "},
}

func TestGetJobInfo(t *testing.T) {
	resource = new(mock.MockResource)

	mux := utils.NewRouter(listRoutes)
	for _, test := range JobInfoTests {

		req, err := http.NewRequest(test.cmd, "http://localhost:8002/"+test.path+"/"+test.querry, nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, test.answercode, w.Code, test.message)
		if w.Code == http.StatusOK {
			assert.Contains(t, w.Body.String(), test.answer, test.message)
		}
	}

}
