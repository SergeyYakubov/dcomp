package server

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"bytes"

	"time"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

var HMACAuthtests = []struct {
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

func writeAuthResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	auth := context.Get(r, "authorizationResponce").(*AuthorizationResponce)
	w.Write([]byte(auth.UserName))
}

func TestProcessHMACAuth(t *testing.T) {

	for _, test := range HMACAuthtests {
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

func TestGenerateJWTToken(t *testing.T) {

	a := NewJWTAuth("hi", "testuser", 0)
	token, _ := a.GenerateToken(nil)
	assert.Equal(t, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ0ZXN0dXNlciJ9.uykaFwU8xMa2"+
		"nXHkgniU5mHls9UvUk6njzzQ1DW5ekg", token, "jwt token")

}

var HJWTAuthtests = []struct {
	Key        string
	User       string
	Duration   time.Duration
	Answercode int
	Message    string
}{
	{"hi", "yakubov", time.Hour, http.StatusOK, "correct auth"},
	{"hi", "yakubov", time.Microsecond, http.StatusUnauthorized, "token expired"},
	{"hih", "yakubov", 1, http.StatusUnauthorized, "wrong key"},
	{"", "yakubov", 1, http.StatusUnauthorized, "auth no header"},
}

func TestProcessJWTAuth(t *testing.T) {
	for _, test := range HJWTAuthtests {
		req, _ := http.NewRequest("POST", "http://blabla", nil)
		a := NewJWTAuth(test.Key, test.User, test.Duration)
		token, _ := a.GenerateToken(req)
		if test.Key != "" {
			req.Header.Add("Authorization", token)
		}
		w := httptest.NewRecorder()
		time.Sleep(time.Second)
		ProcessJWTAuth(http.HandlerFunc(writeAuthResponse), "hi")(w, req)
		assert.Equal(t, test.Answercode, w.Code, test.Message)
		if w.Code == http.StatusOK {
			assert.Contains(t, w.Body.String(), test.User, test.Message)
		}
	}
}

func TestGenerateExternalToken(t *testing.T) {

	a := NewExternalAuth("hi")
	token, _ := a.GenerateToken(nil)
	assert.Equal(t, "hi", token, "external token")

}
