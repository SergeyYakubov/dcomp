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

	"io/ioutil"

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
	{"test.txt", "5", "testuser", "578359205e935a20adb39a18", http.StatusCreated, "receive file"},
	{`blabla/test.txt`, "5", "testuser", "578359205e935a20adb39a18", http.StatusCreated, "receive file with dirname"},
	{`blabla/test.txt`, "5", "testuser", "578359205e935a20adb39a19", http.StatusUnauthorized, "wrong job id"},
	{`blabla/test.txt`, "5", "wronguser", "578359205e935a20adb39a18", http.StatusUnauthorized, "receive file with dirname"},
	{"", "6", "testuser", "578359205e935a20adb39a18", http.StatusBadRequest, "empty filename"},
}

func TestReceiveFiles(t *testing.T) {

	configFileName := `/etc/dcomp/plugins/local/dmd.yaml`
	setDaemonConfiguration(configFileName)
	mux := utils.NewRouter(listRoutes)

	for _, test := range receiveFilesTests {
		m := new(bytes.Buffer)

		var mode os.FileMode = 0600
		binary.Write(m, binary.LittleEndian, &mode)

		b := new(bytes.Buffer)
		b.Write([]byte("Hello"))

		req, err := http.NewRequest("POST", "http://localhost:8002/jobfile/578359205e935a20adb39a18/", b)

		cd := "attachment; filename=" + url.QueryEscape(test.filename)

		req.Header.Set("Content-Disposition", cd)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Content-Mode", url.QueryEscape(m.String()))

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

		if w.Code == http.StatusCreated {
			fname := settings.Resource.BaseDir + "/" + test.jobID + "/" + test.filename
			f, err := os.Open(fname)
			assert.Nil(t, err, "Open file - Should not be error")
			b, err := ioutil.ReadAll(f)
			assert.Contains(t, string(b), "Hello")

			info, err := os.Stat(fname)
			assert.Nil(t, err, "stat file - Should not be error")
			var mode os.FileMode = 0666

			assert.Equal(t, mode, info.Mode())
			os.RemoveAll(settings.Resource.BaseDir + "/" + test.jobID)

		}

	}

}
