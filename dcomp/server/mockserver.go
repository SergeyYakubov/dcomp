// +build !release

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func MockFuncOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, `{"ImageName":"ddd","Script":"aaa","NCPUs":1,"Id":"1","Status":1}`)
}

func MockFuncSimpleString(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Hello")
}

func MockFuncBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "MockFuncBadRequest: bad request", http.StatusBadRequest)
}

func CreateMockServer(srv *Server, mode string) *httptest.Server {
	var ts *httptest.Server
	switch mode {
	case "badreq":
		ts = httptest.NewServer(http.HandlerFunc(MockFuncBadRequest))
	case "string":
		ts = httptest.NewServer(http.HandlerFunc(MockFuncSimpleString))
	default:
		ts = httptest.NewServer(http.HandlerFunc(MockFuncOk))
	}
	srv.ParseUrl(ts.URL)
	return ts
}
