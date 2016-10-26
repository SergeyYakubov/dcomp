package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeDeleteJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	err := resource.DeleteJob(jobID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}
