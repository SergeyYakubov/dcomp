package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	b, err := dbServer.CommandDelete("jobs" + "/" + jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Write(b.Bytes())
	}
}
