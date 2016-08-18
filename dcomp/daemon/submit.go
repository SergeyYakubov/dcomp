package daemon

import (
	"net/http"

	"encoding/json"

	"errors"

	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func routeSubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobDescription

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	job, err := submitJob(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&job)

}

func submitJob(t structs.JobDescription) (job structs.JobInfo, err error) {

	job, err = addJobToDatabase(t)
	if err != nil {
		return
	}

	prio, err := findResources(t)
	if err != nil {
		return
	}

	job.Resource, err = submitToResource(job, prio.Sort())
	if err != nil {
		return
	}

	job.Status = 1

	return
}

func addJobToDatabase(t structs.JobDescription) (job structs.JobInfo, err error) {
	b, err := dbServer.CommandPost("jobs", &t)

	if err != nil {
		return
	}

	err = json.NewDecoder(b).Decode(&job)
	return
}

func submitToResource(job structs.JobInfo, prio []string) (res string, err error) {
	for i := range prio {
		r, ok := resources[prio[i]]
		if ok {
			_, e := r.Server.CommandPost("jobs", job)
			if e == nil {
				res = prio[i]
				return
			}
		}
	}

	err = errors.New("no resource available")
	return
}

func findResources(t structs.JobDescription) (prio structs.ResourcePrio, err error) {

	b, err := estimatorServer.CommandPost("estimations", &t)
	if err != nil {
		return
	}
	err = json.NewDecoder(b).Decode(&prio)
	return
}
