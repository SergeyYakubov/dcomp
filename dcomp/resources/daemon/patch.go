package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func routePatchJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	var patch structs.PatchJob
	if err := patch.Decode(r.Body); err != nil {
		http.Error(w, "routePatchJob: cannot decode patch", http.StatusBadRequest)
		return
	}

	if err := resource.PatchJob(jobID, patch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}
