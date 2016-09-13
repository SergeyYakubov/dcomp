// +build !release

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	//	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/utils"
)

var listRoutes = utils.Routes{
	utils.Route{
		"GetAllJobs",
		"GET",
		"/jobs/",
		MockFuncGetAll,
	},
	utils.Route{
		"GetJob",
		"GET",
		"/jobs/{jobID}/",
		MockFuncGet,
	},
	utils.Route{
		"PatchJob",
		"PATCH",
		"/jobs/{jobID}/",
		MockFuncGet,
	},

	utils.Route{
		"DeleteJob",
		"DELETE",
		"/jobs/{jobID}/",
		MockFuncDelete,
	},
	utils.Route{
		"SubmitJob",
		"POST",
		"/jobs/",
		MockFuncSubmit,
	},
	utils.Route{
		"EstimateJob",
		"POST",
		"/estimations/",
		MockFuncEstimate,
	},
}

func MockFuncSubmit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, `{"ImageName":"submittedimage","Script":"aaa","NCPUs":1,"Id":"578359205e935a20adb39a18","Status":1}`)
}

func MockFuncEstimate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"Local":100,"HPC":0,"Cloud":10,"Batch":0}`)
}

func MockFuncGetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	showFinished := r.URL.Query().Get("finished")
	if showFinished == "true" {
		fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a19"}]`)
	} else {
		fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a18"}]`)
	}
}

func MockFuncGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]

	if jobID == "678359205e935a20adb39a18" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"Status":102}`)
		return
	} else if jobID != "578359205e935a20adb39a18" {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `[{"Id":"578359205e935a20adb39a18"}]`)
}

func MockFuncDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	if jobID != "578359205e935a20adb39a18" {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateMockServer(srv *Server) *httptest.Server {
	var ts *httptest.Server
	mux := utils.NewRouter(listRoutes)
	ts = httptest.NewServer(http.HandlerFunc(mux.ServeHTTP))
	srv.parseUrl(ts.URL)
	return ts
}
