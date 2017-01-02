package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"net/http"

	"net/http/httptest"

	"time"

	"os/user"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"net/url"
	"os"
)

type sendFileNamesRequest struct {
	name       string
	user       string
	recursive  bool
	jobID      string
	answercode int
	result     string
	message    string
}

var sendFilesTests = []sendFileNamesRequest{
	{"test.txt", "testuser", false, "578359205e935a20adb39a18", http.StatusOK, "test.txt", "receive single file"},
	{"/", "testuser", true, "578359205e935a20adb39a18", http.StatusOK, "test.sh", "receive all files,recursive"},
	{url.QueryEscape("/"), "testuser", true, "578359205e935a20adb39a18", http.StatusOK, "test.sh", "receive all files,recursive"},
	{"/", "testuser", false, "578359205e935a20adb39a18", http.StatusOK, "test.txt", "receive all files,recursive"},
	{"test.txt", "wronguser", false, "578359205e935a20adb39a18", http.StatusUnauthorized, "authorization error", "receive single file, wrong user"},
	{"test.txt", "testuser", true, "578359205e935a20adb39a18", http.StatusOK, "test.txt", "receive single file, with recusrion"},
	{"test.tx", "testuser", false, "578359205e935a20adb39a18", http.StatusNotFound, "test.tx", "folder not found"},
	{"test", "testuser", false, "578359205e935a20adb39a18", http.StatusOK, "test", "receive folder"},
	{"test", "testuser", true, "578359205e935a20adb39a18", http.StatusOK, "test.sh", "receive folder, with recusrion"},
	{"tes", "testuser", true, "578359205e935a20adb39a18", http.StatusNotFound, "no such", "folder not found"},
}

func TestSendFiles(t *testing.T) {

	configFileName := `/etc/dcomp/plugins/local/local_dmd.yaml`
	setDaemonConfiguration(configFileName)
	mux := utils.NewRouter(listRoutes)

	path := settings.Resource.BaseDir + `/578359205e935a20adb39a18/test`
	os.MkdirAll(path, 0777)
	os.Create(path + `/test.sh`)
	os.Create(settings.Resource.BaseDir + `/578359205e935a20adb39a18` + `/test.txt`)

	for _, test := range sendFilesTests {

		cmdstr := "http://localhost:8002/jobfile/578359205e935a20adb39a18"
		cmdstr += "/?path=" + test.name + "&nameonly=true"
		if test.recursive {
			cmdstr += "&recursive=true"
		}

		req, err := http.NewRequest("GET", cmdstr, nil)

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
		assert.Contains(t, w.Body.String(), test.result)
		assert.NotContains(t, w.Body.String(), settings.Resource.BaseDir+`/578359205e935a20adb39a18/`)

	}

	os.RemoveAll(settings.Resource.BaseDir + `/578359205e935a20adb39a18`)
}
