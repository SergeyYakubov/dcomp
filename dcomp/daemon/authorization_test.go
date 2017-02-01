package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
)

var userAuthtests = []struct {
	Token      string
	Answercode int
	Answer     string
	Message    string
}{
	{"Basic wronguser", 401, "not", "user not allowed"},
	{"Basic user", 200, "user", "correct auth"},
	{"Wrong test", 401, "type", "wrong auth type"},
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	resp := r.Context().Value("authorizationResponce").(*server.AuthorizationResponce)
	fmt.Fprintln(w, resp.UserName)
}

func TestProcessUserAuth(t *testing.T) {
	var srv server.Server
	ts := server.CreateMockServer(&srv)
	defer ts.Close()
	authServer = srv
	for _, test := range userAuthtests {
		req, _ := http.NewRequest("POST", "http://blabla", nil)
		token := test.Token
		if token != "" {
			req.Header.Add("Authorization", token)
		}
		w := httptest.NewRecorder()
		ProcessUserAuth(http.HandlerFunc(ok))(w, req)
		assert.Equal(t, test.Answercode, w.Code, test.Message)
		assert.Contains(t, w.Body.String(), test.Answer, test.Message)
	}
}
