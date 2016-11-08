package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"net/http"
	"os/user"
)

type Auth interface {
	GenerateToken(*http.Request) (string, error)
}

type BasicAuth struct{}

func (b BasicAuth) GenerateToken(r *http.Request) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

type HMACAuth struct {
	Key string
}

func NewBasicAuth() *BasicAuth {
	a := BasicAuth{}
	return &a
}

func NewHMACAuth(key string) *HMACAuth {
	a := HMACAuth{key}
	return &a
}

func (h HMACAuth) GenerateToken(r *http.Request) (string, error) {

	u := utils.StripURL(r.URL)
	mac := hmac.New(sha256.New, []byte(h.Key))
	mac.Write([]byte(u))

	sha := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return "HMAC-SHA-256 " + sha, nil
}
