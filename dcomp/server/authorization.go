package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/user"
	"strings"
	"time"

	"github.com/gorilla/context"

	"github.com/dgrijalva/jwt-go"
)

type AuthorizationRequest struct {
	Token   string
	Command string
	URL     string
}

type AuthorizationResponce struct {
	UserName string
}

type Auth interface {
	GenerateToken(*http.Request) (string, error)
}

type BasicAuth struct {
	forcedUsername string
}

type ExternalAuth struct {
	token string
}

func NewExternalAuth(token string) *ExternalAuth {
	a := ExternalAuth{}
	a.token = token
	return &a
}

func (b ExternalAuth) GenerateToken(r *http.Request) (string, error) {
	return b.token, nil
}

func (b BasicAuth) GenerateToken(r *http.Request) (string, error) {
	if b.forcedUsername == "" {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		return "Basic " + user.Username, nil
	} else {
		return "Basic " + b.forcedUsername, nil
	}

}

type HMACAuth struct {
	Key string
}

func NewBasicAuth(fn ...string) *BasicAuth {
	a := BasicAuth{}
	if len(fn) > 0 {
		a.forcedUsername = fn[0]
	}
	return &a
}

func NewHMACAuth(key string) *HMACAuth {
	a := HMACAuth{key}
	return &a
}

func generateHMACToken(r *http.Request, key string) string {
	u := stripURL(r.URL)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(u))

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		mac.Write(bodyBytes)
	}

	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func (h HMACAuth) GenerateToken(r *http.Request) (string, error) {

	sha := generateHMACToken(r, h.Key)
	return "HMAC-SHA-256 " + sha, nil
}

func stripURL(u *url.URL) string {
	s := u.Path + u.RawQuery
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, "?", "", -1)
	return s

}

func SplitAuthToken(s string) (authType, token string, err error) {
	keys := strings.Split(s, " ")

	if len(keys) != 2 {
		err = errors.New("authorization error - wrong token")
		return
	}

	authType = keys[0]
	token = keys[1]
	return
}

func ExtractAuthInfo(r *http.Request) (authType, token string, err error) {
	t := r.Header.Get("Authorization")

	if t == "" {
		err = errors.New("authorization error - empty auth header")
		return
	}

	return SplitAuthToken(t)
}

func checkHMACToken(r *http.Request, token, key string) bool {

	if token == "" {
		return false
	}

	generated_token := generateHMACToken(r, key)
	return token == generated_token
}

func ProcessHMACAuth(fn http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authType, token, err := ExtractAuthInfo(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if authType == "HMAC-SHA-256" {
			if !checkHMACToken(r, token, key) {
				http.Error(w, "Internal authorization error - tocken does not match", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Internal authorization error - wrong auth type", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

type JWTAuth struct {
	Key      string
	User     string
	Duration time.Duration
}

func NewJWTAuth(key, user string, d time.Duration) *JWTAuth {
	a := JWTAuth{key, user, d}
	return &a
}

func (t JWTAuth) GenerateToken(r *http.Request) (string, error) {

	claims := &jwt.StandardClaims{
		Id: t.User,
	}

	if t.Duration > 0 {
		claims.ExpiresAt = time.Now().Add(t.Duration).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.Key))

	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}

func ProcessJWTAuth(fn http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authType, token, err := ExtractAuthInfo(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if authType == "Bearer" {
			if resp, ok := checkJWTToken(r, token, key); !ok {
				http.Error(w, "Internal authorization error - tocken does not match", http.StatusUnauthorized)
				return
			} else {
				context.Set(r, "authorizationResponce", &resp)
				defer context.Clear(r)
			}
		} else {
			http.Error(w, "Internal authorization error - wrong auth type", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func checkJWTToken(r *http.Request, token, key string) (AuthorizationResponce, bool) {

	var resp AuthorizationResponce
	if token == "" {
		return resp, false
	}

	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err == nil && t.Valid {
		resp.UserName = t.Claims.(*jwt.StandardClaims).Id
		return resp, true
	}

	return resp, false
}
