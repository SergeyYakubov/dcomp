package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"

	"encoding/base64"

	"strings"

	"github.com/gorilla/mux"
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

func Auth(fn http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sha := r.Header.Get("Authorization")
		reqToken, err := base64.URLEncoding.DecodeString(sha)
		if err != nil || sha == "" {
			http.Error(w, "Internal authorization error - empty tocken", http.StatusUnauthorized)
			return
		}
		message := r.URL.Path + r.URL.RawQuery
		message = strings.Replace(message, "/", "", -1)
		message = strings.Replace(message, "?", "", -1)

		ok := checkMAC([]byte(message), reqToken, []byte(key))
		if !ok {
			http.Error(w, "Internal authorization error - tocken does not match", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
