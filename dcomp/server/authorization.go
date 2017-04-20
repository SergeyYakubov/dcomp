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

	"context"

	"github.com/sergeyyakubov/dcomp/dcomp/utils"

	"github.com/dgrijalva/jwt-go"
)

type AuthorizationRequest struct {
	Token   string
	Command string
	URL     string
}

type AuthorizationResponce struct {
	Status       int
	StatusText   string
	UserName     string
	Token        string
	ValidityTime int
}

type Auth interface {
	GenerateToken(...interface{}) (string, error)
	Name() string
}

func (a *BasicAuth) Name() string {
	return "Basic"
}

func (a *ExternalAuth) Name() string {
	return "External"
}

func (a *HMACAuth) Name() string {
	return "HMAC-SHA-256"
}

func (a *GSSAPIAuth) Name() string {
	return "Negotiate"
}

func (a *JWTAuth) Name() string {
	return "Bearer"
}

type BasicAuth struct {
	forcedUsername string
}

type NoAuth struct {
}

func (a *NoAuth) Name() string {
	return "None"
}
func (a *NoAuth) GenerateToken(...interface{}) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return "None " + user.Username, nil
}

func NewNoAuth() *NoAuth {
	a := NoAuth{}
	return &a
}

type ExternalAuth struct {
	Token string
}

func NewExternalAuth(token string) *ExternalAuth {
	a := ExternalAuth{}
	a.Token = token
	return &a
}

func (b ExternalAuth) GenerateToken(...interface{}) (string, error) {
	return b.Token, nil
}

func (b BasicAuth) GenerateToken(...interface{}) (string, error) {
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

func (h HMACAuth) GenerateToken(val ...interface{}) (string, error) {
	if len(val) != 1 {
		return "", errors.New("Wrong claims")
	}
	claims, ok := val[0].(*CustomClaims)
	if !ok {
		return "", errors.New("Wrong claims")
	}

	r := claims.ExtraClaims.(*http.Request)
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

	if t != "" {
		return SplitAuthToken(t)
	}

	cookie, err := r.Cookie("Authorization")
	if err == nil {
		return SplitAuthToken(cookie.Value)
	}

	err = errors.New("no authorization info")
	return

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

type CustomClaims struct {
	jwt.StandardClaims
	Duration    time.Duration
	ExtraClaims interface{}
}

type JobClaim struct {
	AuthorizationResponce
	JobInd string
}

type JWTAuth struct {
	Key string
}

func NewJWTAuth(key string) *JWTAuth {
	a := JWTAuth{key}
	return &a
}

func (t JWTAuth) GenerateToken(val ...interface{}) (string, error) {
	if len(val) != 1 {
		return "", errors.New("Wrong claims")
	}
	claims, ok := val[0].(*CustomClaims)
	if !ok {
		return "", errors.New("Wrong claims")
	}

	if claims.Duration > 0 {
		claims.ExpiresAt = time.Now().Add(claims.Duration).Unix()
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

		ctx := r.Context()

		if authType == "Bearer" {
			if claims, ok := CheckJWTToken(token, key); !ok {
				http.Error(w, "Internal authorization error - tocken does not match", http.StatusUnauthorized)
				return
			} else {
				ctx = context.WithValue(ctx, "JobClaim", claims)
			}
		} else {
			http.Error(w, "Internal authorization error - wrong auth type", http.StatusUnauthorized)
			return
		}
		fn(w, r.WithContext(ctx))
	}
}

func CheckJWTToken(token, key string) (jwt.Claims, bool) {

	if token == "" {
		return nil, false
	}

	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err == nil && t.Valid {
		return t.Claims, true
	}

	return nil, false
}

func JobClaimFromContext(r *http.Request, val interface{}) error {
	c := r.Context().Value("JobClaim")

	if c == nil {
		return errors.New("Empty context")
	}

	claim := c.(*CustomClaims)

	return utils.MapToStruct(claim.ExtraClaims.(map[string]interface{}), val)
}

type GSSAPIAuth struct {
}

func (b GSSAPIAuth) GenerateToken(...interface{}) (string, error) {

	data, err := utils.GetGSSAPIToken("dcomp")

	if err != nil {
		return "", err
	}

	token := "Negotiate" + " " + base64.StdEncoding.EncodeToString(data)
	return token, nil
}

func NewGSSAPIAuth() *GSSAPIAuth {
	a := GSSAPIAuth{}
	return &a
}
