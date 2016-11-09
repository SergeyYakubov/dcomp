package utils

import (
	"net/http"

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
