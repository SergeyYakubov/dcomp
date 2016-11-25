package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"net/http/httptest"

	"bytes"

	"net/url"

	"encoding/binary"
	"os"
	"time"

	"os/user"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type receiveFilesRequest struct {
	filename   string
	filelength string
	user       string
	jobID      string
	answercode int
	message    string
}

var receiveFilesTests = []receiveFilesRequest{
	{"test.txt", "5", "testuser", "578359205e935a20adb39a18", http.StatusOK, "receive file"},
	{`blabla/test.txt`, "5", "testuser", "578359205e935a20adb39a18", http.StatusOK, "receive file with dirname"},
	{`blabla/test.txt`, "5", "testuser", "578359205e935a20adb39a19", http.StatusUnauthorized, "wrong job id"},
	{`blabla/test.txt`, "5", "wronguser", "578359205e935a20adb39a18", http.StatusUnauthorized, "receive file with dirname"},
	{"test.txt", "6", "testuser", "578359205e935a20adb39a18", http.StatusInternalServerError, "wrong length"},
	{"", "6", "testuser", "578359205e935a20adb39a18", http.StatusBadRequest, "empty filename"},
}

func TestReceiveFiles(t *testing.T) {

	configFileName := `/etc/dcomp/plugins/local/local_dmd.yaml`
	setDaemonConfiguration(configFileName)
	mux := utils.NewRouter(listRoutes)

	for _, test := range receiveFilesTests {
		b := new(bytes.Buffer)

		var mode os.FileMode = 0600
		binary.Write(b, binary.LittleEndian, &mode)

		b.Write([]byte("Hello"))

		req, err := http.NewRequest("POST", "http://localhost:8002/jobfile/578359205e935a20adb39a18/", b)

		cd := "attachment; filename=" + url.QueryEscape(test.filename)

		req.Header.Set("Content-Disposition", cd)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", test.filelength)

		assert.Nil(t, err, "Should not be error")

		if test.user == "testuser" {
			u, _ := user.Current()
			test.user = u.Username
		}

		var claim server.JobClaim
		claim.UserName = test.user
		claim.JobInd = test.jobID

		var c server.CustomClaims
		c.ExtraClaims = &claim
		c.Duration = time.Hour

		auth := server.NewJWTAuth(settings.Daemon.Key)
		token, _ := auth.GenerateToken(&c)

		if test.user != "wronguser" {
			req.Header.Add("Authorization", token)
		}

		w := httptest.NewRecorder()
		f := server.ProcessJWTAuth(mux.ServeHTTP, settings.Daemon.Key)
		f(w, req)

		assert.Equal(t, test.answercode, w.Code, test.message)
	}

}
