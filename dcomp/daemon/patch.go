package daemon

import (
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func routePatchJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res, ok := resources[job.Resource]
	if !ok {
		http.Error(w, "Resource unvailable "+job.Resource, http.StatusInternalServerError)
	}

	var patch structs.PatchJob
	if err := patch.Decode(r.Body); err != nil {
		http.Error(w, "routePatchJob: cannot decode patch", http.StatusBadRequest)
		return
	}

	_, _, err = res.Server.CommandPatch("jobs"+"/"+job.Id, &patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
