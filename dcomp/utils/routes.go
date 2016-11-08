package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"

	"encoding/base64"

	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/url"
)

type Routes []Route

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewRouter(listRoutes Routes) *mux.Router {
	router := mux.NewRouter()
	for _, route := range listRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func StripURL(u *url.URL) string {
	s := u.Path + u.RawQuery
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, "?", "", -1)
	return s

}

func extractAuthorizationInfo(r *http.Request) (authType, token string, err error) {
	t := r.Header.Get("Authorization")

	keys := strings.Split(t, " ")

	if len(keys) != 2 {
		err = errors.New("Internal authorization error - wrong token")
		return
	}

	authType = keys[0]
	token = keys[1]
	return
}

func checkHMACToken(r *http.Request, token, key string) bool {

	reqToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil || token == "" {
		return false
	}
	message := StripURL(r.URL)

	return checkMAC([]byte(message), reqToken, []byte(key))
}

func HMACAuth(fn http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authType, token, err := extractAuthorizationInfo(r)

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
