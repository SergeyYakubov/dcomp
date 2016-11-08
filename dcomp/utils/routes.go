package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"

	"encoding/base64"

	"strings"

	"github.com/gorilla/mux"
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

func Auth(fn http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sha := r.Header.Get("Authorization")
		reqToken, err := base64.URLEncoding.DecodeString(sha)
		if err != nil || sha == "" {
			http.Error(w, "Internal authorization error - empty tocken", http.StatusUnauthorized)
			return
		}
		message := StripURL(r.URL)

		ok := checkMAC([]byte(message), reqToken, []byte(key))
		if !ok {
			http.Error(w, "Internal authorization error - tocken does not match", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
