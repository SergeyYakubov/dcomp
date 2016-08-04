package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func routeEstimateJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobDescription

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	prio := estimate(t)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(prio)
	w.Write(b.Bytes())
}

// estimate uses very simple algorithm to assign priorities based on CPU number required for the job
func estimate(job structs.JobDescription) (prio structs.ResourcePrio) {
	prio = make(structs.ResourcePrio)
	switch {
	case job.NCPUs == 1:
		prio["HPC"] = 1
		prio["Batch"] = 10
		prio["Cloud"] = 0
	case job.NCPUs <= 8:
		prio["HPC"] = 5
		prio["Batch"] = 5
		prio["Cloud"] = 0
	default:
		prio["HPC"] = 10
		prio["Batch"] = 0
		prio["Cloud"] = 0
	}
	return prio
}
