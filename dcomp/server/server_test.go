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

	var srv Srv

	srv.Port = -4
	b, err := srv.PostCommand("", nil)
	assert.NotNil(t, err, "Should be error in http.post")

	ts := CreateMockServer(&srv)
	defer ts.Close()

	b, err = srv.PostCommand("", ts)
	assert.NotNil(t, err, "Should be error in json encoder")

	// nil is actually a bad option but since we use mock server we cannot check it
	b, err = srv.PostCommand("", nil)
	assert.Equal(t, "{\"ImageName\":\"ddd\",\"Script\":\"aaa\",\"NCPUs\":1,\"Id\":1,\"Status\":1}\n",
		b.String(), "")

	srv.Port = 10000
	b, err = srv.PostCommand("", nil)
	assert.Contains(t, err.Error(), "connection refused", "")

	srv.Host = "aaa"
	b, err = srv.PostCommand("", nil)
	assert.Contains(t, err.Error(), "no such host", "")

}
