package server

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"bytes"

	"github.com/stretchr/testify/assert"
)

var hmacAuthtests = []struct {
	Key        string
	Answercode int
	Message    string
}{
	{"hi", http.StatusOK, "correct auth"},
	{"hih", http.StatusUnauthorized, "wrong key"},
	{"", http.StatusUnauthorized, "auth no header"},
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestProcessHMACAuth(t *testing.T) {

	for _, test := range hmacAuthtests {
		req, _ := http.NewRequest("POST", "http://blabla", nil)
		a := NewHMACAuth(test.Key)
		token, _ := a.GenerateToken(req)
		if test.Key != "" {
			req.Header.Add("Authorization", token)
		}
		w := httptest.NewRecorder()
		ProcessHMACAuth(http.HandlerFunc(ok), "hi")(w, req)
		assert.Equal(t, test.Answercode, w.Code, test.Message)
	}
}

func TestGenerateToken(t *testing.T) {
	buf := bytes.NewBuffer([]byte("aaa"))
	req, _ := http.NewRequest("POST", "http://blabla", buf)
	a := NewHMACAuth("hi")

	token, _ := a.GenerateToken(req)
	assert.Equal(t, "HMAC-SHA-256 SXNIGkudNaJZdsY4zVCjFXwz7laxitp-ZsUgrvd5Acc=", token, "hmac token")

	b := NewBasicAuth()
	token, _ = b.GenerateToken(req)

	assert.Contains(t, token, "Basic", "basic token")

	b = NewBasicAuth("test")
	token, _ = b.GenerateToken(req)

	assert.Equal(t, token, "Basic test", "basic token")
}
