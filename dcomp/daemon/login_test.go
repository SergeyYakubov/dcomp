package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

var loginTests = []struct {
	Token      string
	Answercode int
	Answer     string
	Message    string
}{
	{"Basic wronguser", 401, "not", "user not allowed"},
	{"Basic user", 200, "blabla", "correct auth"},
	{"Wrong test", 400, "wrong", "wrong auth type"},
}

func TestLogin(t *testing.T) {
	mux := utils.NewRouter(listRoutes)

	//setConfiguration()

	var srv server.Server
	ts := server.CreateMockServer(&srv)
	defer ts.Close()
	authServer = srv
	authServer.SetAuth(server.NewBasicAuth())

	for _, test := range loginTests {

		req, _ := http.NewRequest("GET", "http://localhost:8000/login/", nil)

		token := test.Token
		if token != "" {
			req.Header.Add("Authorization", token)
		}

		w := httptest.NewRecorder()
		f := ProcessUserAuth(mux.ServeHTTP)
		f(w, req)
		assert.Equal(t, test.Answercode, w.Code, test.Message)
		assert.Contains(t, w.Body.String(), test.Answer, test.Message)
	}
}
