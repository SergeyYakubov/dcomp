package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
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
	prio["cloud"] = 0
	prio["batch"] = 0
	prio["maxwell"] = 0
	prio["local"] = 0
	switch {
	case job.NCPUs == 1:
		prio["maxwell"] = 1
		prio["batch"] = 10
		prio["local"] = 20
	case job.NCPUs <= 8:
		prio["maxwell"] = 5
		prio["batch"] = 5
		prio["local"] = 20
	}
	if job.Resource != "" {
		prio[job.Resource] = 100
	}

	return prio

}
