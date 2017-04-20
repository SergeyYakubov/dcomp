package server

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"bytes"

	"time"

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
	var jc JobClaim
	JobClaimFromContext(r, &jc)
	w.Write([]byte(jc.UserName))
	w.Write([]byte(jc.JobInd))
}

func TestProcessHMACAuth(t *testing.T) {

	for _, test := range HMACAuthtests {
		req, _ := http.NewRequest("POST", "http://blabla", nil)
		a := NewHMACAuth(test.Key)
		token, _ := a.GenerateToken(&CustomClaims{ExtraClaims: req})
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

	token, _ := a.GenerateToken(&CustomClaims{ExtraClaims: req})
	assert.Equal(t, "HMAC-SHA-256 SXNIGkudNaJZdsY4zVCjFXwz7laxitp-ZsUgrvd5Acc=", token, "hmac token")

	b := NewBasicAuth()
	token, _ = b.GenerateToken(&CustomClaims{ExtraClaims: req})

	assert.Contains(t, token, "Basic", "basic token")

	b = NewBasicAuth("test")
	token, _ = b.GenerateToken(&CustomClaims{ExtraClaims: req})

	assert.Equal(t, token, "Basic test", "basic token")
}

func TestGenerateJWTToken(t *testing.T) {

	a := NewJWTAuth("hi")
	token, _ := a.GenerateToken((&CustomClaims{Duration: 0, ExtraClaims: nil}))
	assert.Equal(t, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEdXJhdGlvbiI"+
		"6MCwiRXh0cmFDbGFpbXMiOm51bGx9.JJcqNZciIDILk-A2sJZCY1sND458bcjNv6tXC2jxric",
		token, "jwt token")

}

var HJWTAuthtests = []struct {
	Mode       string
	Key        string
	User       string
	jobID      string
	Duration   time.Duration
	Answercode int
	Message    string
}{
	{"header", "hi", "testuser", "123", time.Hour, http.StatusOK, "correct auth - header"},
	{"cookie", "hi", "testuser", "123", time.Hour, http.StatusOK, "correct auth - cookie"},
	{"header", "hi", "testuser", "123", time.Microsecond, http.StatusUnauthorized, "token expired"},
	{"header", "hih", "testuser", "123", 1, http.StatusUnauthorized, "wrong key"},
	{"", "hi", "testuser", "123", 1, http.StatusUnauthorized, "auth no header"},
}

func TestProcessJWTAuth(t *testing.T) {
	for _, test := range HJWTAuthtests {
		req, _ := http.NewRequest("POST", "http://blabla", nil)

		var claim JobClaim
		claim.UserName = test.User
		claim.JobInd = test.jobID

		a := NewJWTAuth(test.Key)

		token, _ := a.GenerateToken((&CustomClaims{Duration: test.Duration, ExtraClaims: &claim}))
		if test.Mode == "header" {
			req.Header.Add("Authorization", token)
		}

		if test.Mode == "cookie" {
			c := http.Cookie{Name: "Authorization", Value: token}
			req.AddCookie(&c)
		}

		w := httptest.NewRecorder()
		if test.Duration == time.Microsecond {
			if testing.Short() {
				continue
			}
			time.Sleep(time.Second)
		}
		ProcessJWTAuth(http.HandlerFunc(writeAuthResponse), "hi")(w, req)
		assert.Equal(t, test.Answercode, w.Code, test.Message)
		if w.Code == http.StatusOK {
			assert.Contains(t, w.Body.String(), test.User, test.Message)
			assert.Contains(t, w.Body.String(), test.jobID, test.Message)
		}
	}
}

func TestGenerateExternalToken(t *testing.T) {

	a := NewExternalAuth("hi")
	token, _ := a.GenerateToken()
	assert.Equal(t, "hi", token, "external token")

}
