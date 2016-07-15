package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	b, err := DBServer.CommandGet("jobs")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write(b.Bytes())
	}

}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	b, err := DBServer.CommandGet("jobs" + "/" + jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Write(b.Bytes())
	}

}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	b, err := DBServer.CommandDelete("jobs" + "/" + jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Write(b.Bytes())
	}
}
