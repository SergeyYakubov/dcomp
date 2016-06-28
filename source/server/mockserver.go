package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"stash.desy.de/scm/dc/daemon"
)

func MockFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func MockDaemonRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range daemon.ListRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(http.HandlerFunc(MockFunc))
	}
	return router
}

func CreateMockServer(server *Srv, mode string) *httptest.Server {
	var mux *mux.Router
	switch mode {
	case "daemon":
		mux = MockDaemonRouter()
	}

	ts := httptest.NewServer(http.HandlerFunc(mux.ServeHTTP))
	server.ParseUrl(ts.URL)
	return ts
}
