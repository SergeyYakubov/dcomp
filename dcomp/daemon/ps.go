package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeGetAllJobs(w http.ResponseWriter, r *http.Request) {
	b, err := dbServer.CommandGet("jobs")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write(b.Bytes())
	}

}

func routeGetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	b, err := dbServer.CommandGet("jobs" + "/" + jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Write(b.Bytes())
	}

}
