package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var urltests = []struct {
	Server Server
	Path   string
	Url    string
}{
	{Server{"localhost", 8000, "key"}, "test", "http://localhost:8000/test/"},
	{Server{"localhost", 8000, "key"}, "/test/", "http://localhost:8000/test/"},
	{Server{"localhost", 8000, "key"}, " test ", "http://localhost:8000/test/"},
}

func TestUrl(t *testing.T) {
	for _, test := range urltests {
		srv := test.Server
		assert.Equal(t, srv.url(test.Path), test.Url, "")
	}

}

func TestPostcommand(t *testing.T) {

	var srv Server

	srv.Key = "1234"
	srv.Port = -4
	b, err := srv.CommandPost("jobs", nil)
	assert.NotNil(t, err, "Should be error in http.post")

	ts := CreateMockServer(&srv)
	defer ts.Close()

	b, err = srv.CommandPost("jobs", ts)
	assert.NotNil(t, err, "Should be error in json encoder")

	// nil is actually a bad option but since we use mock server we cannot check it
	b, err = srv.CommandPost("jobs", nil)
	assert.Equal(t, "{\"ImageName\":\"submittedimage\",\"Script\":\"aaa\",\"NCPUs\":1,\"Id\":\"578359205e935a20adb39a18\",\"Status\":1}\n",
		b.String(), "")

	srv.Key = "123"
	b, err = srv.CommandPost("jobs", nil)
	assert.NotNil(t, err, "authorization should fail")

	srv.Key = "1234"
	srv.Port = 10000
	b, err = srv.CommandPost("jobs", nil)
	assert.Contains(t, err.Error(), "connection refused", "")

	//	srv.Host = "aaa"
	//	b, err = srv.CommandPost("jobs", nil)
	//	assert.Contains(t, err.Error(), "lookup", "")

	b, err = srv.CommandPost("jobs", nil)
	assert.NotNil(t, err, "Should be error in responce")

}

type httpRequest struct {
	path string
	body string
	key  string
	msg  string
}

var getTests = []httpRequest{
	{"jobs/578359205e935a20adb39a18", "578359205e935a20adb39a18", "1234", "get job 1"},
	{"jobs/2", "not found", "1234", "wrong job id"},
	{"job", "not found", "1234", "wrong path"},
}

var rmTests = []httpRequest{
	{"jobs/578359205e935a20adb39a18", "", "1234", "get job 1"},
	{"jobs/2", "not found", "1234", "wrong job id"},
	{"job", "not found", "1234", "wrong path"},
}

func TestGetcommand(t *testing.T) {
	var srv Server
	ts := CreateMockServer(&srv)
	defer ts.Close()
	for _, test := range getTests {
		srv.Key = test.key
		b, err := srv.CommandGet(test.path)
		if err != nil {
			assert.Contains(t, err.Error(), test.body, test.msg)
		} else {
			assert.Contains(t, b.String(), test.body, test.msg)
		}

	}
}
func TestDeletecommand(t *testing.T) {
	var srv Server
	ts := CreateMockServer(&srv)
	defer ts.Close()
	for _, test := range rmTests {
		srv.Key = test.key
		b, err := srv.CommandDelete(test.path)
		if err != nil {
			assert.Contains(t, err.Error(), test.body, test.msg)
		} else {
			assert.Contains(t, b.String(), test.body, test.msg)
		}

	}
}

var patchTests = []httpRequest{
	{"jobs/578359205e935a20adb39a18", "578359205e935a20adb39a18", "1234", "patch job 1"},
	{"jobs/2", "not found", "1234", "wrong job id"},
	{"job", "not found", "1234", "wrong path"},
}

func TestPatchcommand(t *testing.T) {
	var srv Server
	ts := CreateMockServer(&srv)
	defer ts.Close()
	for _, test := range patchTests {
		srv.Key = test.key
		err := srv.CommandPatch(test.path, nil)
		if err != nil {
			assert.Contains(t, err.Error(), test.body, test.msg)
		}

	}
}
