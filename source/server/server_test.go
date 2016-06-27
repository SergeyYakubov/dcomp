package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var urltests = []struct {
	Server Srv
	Path   string
	Url    string
}{
	{Srv{"localhost", 8000}, "test", "http://localhost:8000/test/"},
	{Srv{"localhost", 8000}, "/test/", "http://localhost:8000/test/"},
	{Srv{"localhost", 8000}, " test ", "http://localhost:8000/test/"},
}

func TestUrl(t *testing.T) {
	for _, test := range urltests {
		srv := test.Server
		assert.Equal(t, srv.Url(test.Path), test.Url, "")
	}

}

func TestPostcommand(t *testing.T) {

	var server Srv

	server.Port = -4
	str, err := server.PostCommand("jobs", nil)
	assert.NotNil(t, err, "Should be error in http.post")

	ts := CreateMockServer(&server, "daemon")
	defer ts.Close()

	str, err = server.PostCommand("jobs", ts)
	assert.NotNil(t, err, "Should be error in json encoder")

	// nil is actually a bad option but since we use mock server we cannot check it
	str, err = server.PostCommand("jobs", nil)
	assert.Equal(t, "OK\n", str, "")

	str, err = server.PostCommand("job", nil)
	assert.Equal(t, "404 page not found\n", str, "")
}
