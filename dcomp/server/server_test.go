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
	{Server{"localhost", 8000}, "test", "http://localhost:8000/test/"},
	{Server{"localhost", 8000}, "/test/", "http://localhost:8000/test/"},
	{Server{"localhost", 8000}, " test ", "http://localhost:8000/test/"},
}

func TestUrl(t *testing.T) {
	for _, test := range urltests {
		srv := test.Server
		assert.Equal(t, srv.Url(test.Path), test.Url, "")
	}

}

func TestPostcommand(t *testing.T) {

	var srv Server

	srv.Port = -4
	b, err := srv.CommandPost("jobs", nil)
	assert.NotNil(t, err, "Should be error in http.post")

	ts := CreateMockServer(&srv)
	defer ts.Close()

	b, err = srv.CommandPost("jobs", ts)
	assert.NotNil(t, err, "Should be error in json encoder")

	// nil is actually a bad option but since we use mock server we cannot check it
	b, err = srv.CommandPost("jobs", nil)
	assert.Equal(t, "{\"ImageName\":\"ddd\",\"Script\":\"aaa\",\"NCPUs\":1,\"Id\":\"1\",\"Status\":1}\n",
		b.String(), "")

	srv.Port = 10000
	b, err = srv.CommandPost("jobs", nil)
	assert.Contains(t, err.Error(), "connection refused", "")

	srv.Host = "aaa"
	b, err = srv.CommandPost("jobs", nil)
	assert.Contains(t, err.Error(), "no such host", "")

	b, err = srv.CommandPost("jobs", nil)
	assert.NotNil(t, err, "Should be error in responce")

}

type getrequest struct {
	path   string
	body   string
	errmsg string
}

var getTests = []getrequest{
	{"jobs/578359205e935a20adb39a18", "578359205e935a20adb39a18", "get job 1"},
	{"jobs/2", "not found", "wrong format"},
	{"job", "not found", "wrong path"},
}

func TestGetcommand(t *testing.T) {
	var srv Server
	ts := CreateMockServer(&srv)
	defer ts.Close()
	for _, test := range getTests {
		b, err := srv.CommandGet(test.path)
		if err != nil {
			assert.Contains(t, err.Error(), test.body, test.errmsg)
		} else {
			assert.Contains(t, b.String(), test.body, test.errmsg)
		}

	}
}
