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

	"encoding/json"
	"fmt"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
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

func TestReceiveFilesMount(t *testing.T) {
	type request struct {
		fi         structs.FileCopyInfo
		jobID      string
		answercode int
		message    string
	}

	var tests = []request{
				{structs.FileCopyInfo{SourcePath: "aaa.txt", Source: "578359205e935a20adb39a19", DestPath: "data/bbb"}, "578359205e935a20adb39a18", http.StatusCreated, "receive file"},
		{structs.FileCopyInfo{SourcePath: "aaa/aaa.txt", Source: "578359205e935a20adb39a19", DestPath: "data/bbb"}, "578359205e935a20adb39a18", http.StatusCreated, "receive file"},
				{structs.FileCopyInfo{SourcePath: "aaa", Source: "578359205e935a20adb39a19", DestPath: "data/bbb"}, "578359205e935a20adb39a18", http.StatusCreated, "receive file"},
	}

	configFileName := `/etc/dcomp/plugins/local/dmd.yaml`
	setDaemonConfiguration(configFileName)
	mux := utils.NewRouter(listRoutes)

	os.MkdirAll(settings.Resource.BaseDir+"/578359205e935a20adb39a19/aaa", 0777)
	fname := settings.Resource.BaseDir + "/578359205e935a20adb39a19/aaa.txt"
	f, _ := os.Create(fname)
	f.WriteString("Hello")
	f.Close()
	fname = settings.Resource.BaseDir + "/578359205e935a20adb39a19/aaa/aaa.txt"
	f, _ = os.Create(fname)
	f.WriteString("Hello1")
	f.Close()

	for _, test := range tests {

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(test.fi)

		req, err := http.NewRequest("POST", "http://localhost:8002/jobfile/578359205e935a20adb39a18/?mode=mount", b)

		assert.Nil(t, err, "Should not be error")

		u, _ := user.Current()

		var claim server.JobClaim
		claim.UserName = u.Username
		claim.JobInd = test.jobID

		var c server.CustomClaims
		c.ExtraClaims = &claim
		c.Duration = time.Hour

		auth := server.NewJWTAuth(settings.Daemon.Key)
		token, _ := auth.GenerateToken(&c)

		req.Header.Add("Authorization", token)

		w := httptest.NewRecorder()
		f := server.ProcessJWTAuth(mux.ServeHTTP, settings.Daemon.Key)
		f(w, req)

		assert.Equal(t, test.answercode, w.Code, test.message)
		fmt.Println(w.Body.String())

		if w.Code == http.StatusCreated {
			var fname string
			if test.fi.SourcePath == "aaa" {
				fname = settings.Resource.BaseDir + "/" + test.jobID + "/" + test.fi.DestPath + "/aaa.txt"
			} else {
				fname = settings.Resource.BaseDir + "/" + test.jobID + "/" + test.fi.DestPath
			}

			f, err := os.Open(fname)
			assert.Nil(t, err, "Open file - Should not be error")
			b, err := ioutil.ReadAll(f)
			assert.Contains(t, string(b), "Hello")
		}
		os.RemoveAll(settings.Resource.BaseDir + "/" + test.jobID)
	}

	os.RemoveAll(settings.Resource.BaseDir + "/578359205e935a20adb39a19")

}
