package cli

import (
	"../daemon"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func ParseUrl(s string) {
	u, _ := url.Parse(s)
	host, port, _ := net.SplitHostPort(u.Host)
	Server.host = host
	Server.port, _ = strconv.Atoi(port)

}

func PrepareMockServer() *httptest.Server {
	mux := daemon.NewRouter()
	ts := httptest.NewServer(http.HandlerFunc(mux.ServeHTTP))
	ParseUrl(ts.URL)
	return ts

}

func TestPostcommand(t *testing.T) {
	OutBuf = new(bytes.Buffer)

	Server.port = -4
	err := Server.PostCommand("jobs", nil)
	assert.NotNil(t, err, "Should be error in http.post")

	ts := PrepareMockServer()
	defer ts.Close()

	err = Server.PostCommand("jobs", ts)
	assert.NotNil(t, err, "Should be error in json encoder")

	Server.PostCommand("jobs", nil)
	assert.Equal(t, "Bad request\n", OutBuf.(*bytes.Buffer).String(), "")

	OutBuf.(*bytes.Buffer).Reset()

	Server.PostCommand("job", nil)
	assert.Equal(t, "404 page not found\n", OutBuf.(*bytes.Buffer).String(), "")
}

/*func TestSubmitJob(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(daemon.SubmitJob))
	parseUrl(ts.URL)
	defer ts.Close()
	Server.PostCommand("", nil)
}*/
