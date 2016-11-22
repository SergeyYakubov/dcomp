package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"net/http/httptest"

	"bytes"

	"net/url"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"time"
)

type receiveFilesRequest struct {
	filename   string
	filelength string
	user       string
	answercode int
	message    string
}

var receiveFilesTests = []receiveFilesRequest{
	{"test.txt", "5", "testuser", http.StatusOK, "receive file"},
	{`blabla/test.txt`, "5", "testuser", http.StatusOK, "receive file with dirname"},
	{`blabla/test.txt`, "5", "wronguser", http.StatusUnauthorized, "receive file with dirname"},
	{"test.txt", "6", "", http.StatusInternalServerError, "wrong length"},
	{"", "6", "", http.StatusBadRequest, "empty filename"},
}

func TestReceiveFiles(t *testing.T) {

	mux := utils.NewRouter(listRoutes)

	for _, test := range receiveFilesTests {
		b := new(bytes.Buffer)
		b.Write([]byte("Hello"))

		req, err := http.NewRequest("POST", "http://localhost:8002/jobfile/578359205e935a20adb39a18/", b)

		cd := "attachment; filename=" + url.QueryEscape(test.filename)

		req.Header.Set("Content-Disposition", cd)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", test.filelength)

		assert.Nil(t, err, "Should not be error")

		auth := server.NewJWTAuth(settings.Daemon.Key, "test", time.Hour)
		token, _ := auth.GenerateToken(req)
		if test.user != "wronguser" {
			req.Header.Add("Authorization", token)
		}

		w := httptest.NewRecorder()
		f := server.ProcessJWTAuth(mux.ServeHTTP, settings.Daemon.Key)
		f(w, req)

		assert.Equal(t, test.answercode, w.Code, test.message)
	}

}
