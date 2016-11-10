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
	{server.AuthorizationRequest{Token: "Basic test"}, http.StatusOK, "test", "correct auth"},
	{server.AuthorizationRequest{Token: "Basic"}, http.StatusUnauthorized, "token", "wrong token"},
	{server.AuthorizationRequest{Token: "BlaBla test"}, http.StatusUnauthorized, "type", "wrong token type"},
	{server.AuthorizationRequest{Token: ""}, http.StatusUnauthorized, "token", "empty token"},
	{server.AuthorizationRequest{Token: "nil"}, http.StatusUnauthorized, "bad", "no request body"},
}

func TestAuthorizeJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)

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
		assert.Equal(t, test.answerCode, w.Code, test.message)
		assert.Contains(t, w.Body.String(), test.answerMessage, test.message)
	}
}
