package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]

	if err := db.DeleteRecordByID(jobID); err != nil {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusNotFound)
	}
}
