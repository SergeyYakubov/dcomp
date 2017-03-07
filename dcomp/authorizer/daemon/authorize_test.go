package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"io"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

type request struct {
	req           server.AuthorizationRequest
	answerCode    int
	answerMessage string
	message       string
}

var authorizeTests = []request{
	{server.AuthorizationRequest{Token: "None test"}, http.StatusOK, "test", "basic auth"},
	{server.AuthorizationRequest{Token: "Negotiate test"}, http.StatusBadRequest, "not defined", "kerb auth, not defined context"},
	{server.AuthorizationRequest{Token: "Basic"}, http.StatusUnauthorized, "token", "wrong token"},
	{server.AuthorizationRequest{Token: "BlaBla test"}, http.StatusUnauthorized, "type", "wrong token type"},
	{server.AuthorizationRequest{Token: ""}, http.StatusUnauthorized, "token", "empty token"},
	{server.AuthorizationRequest{Token: "nil"}, http.StatusBadRequest, "bad", "no request body"},
	{server.AuthorizationRequest{Token: ""}, http.StatusUnauthorized, "wrong", "empty token"},
}

func TestAuthorizeJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	tconf := configFile
	configFile = "config_test.yaml"
	setDaemonConfiguration()

	for _, test := range authorizeTests {

		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(test.req); err != nil {
			t.Fail()
		}

		var reader io.Reader = b
		if test.req.Token == "nil" {
			reader = nil
		}

		req, err := http.NewRequest("POST", "http://localhost:8002/authorize/", reader)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			assert.Equal(t, test.answerCode, w.Code, test.message)
			assert.Contains(t, w.Body.String(), test.answerMessage, test.message)

		} else {

			var resp server.AuthorizationResponce
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, test.answerCode, resp.Status, test.message)
			if resp.Status != http.StatusOK {
				assert.Contains(t, resp.StatusText, test.answerMessage, test.message)
			} else {
				assert.Contains(t, resp.UserName, test.answerMessage, test.message)
			}

		}

	}
	configFile = tconf
	setDaemonConfiguration()

}
