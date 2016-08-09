package daemon

import (
	"net/http"

	"encoding/json"

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

func addJobToDatabase(t structs.JobDescription) (job structs.JobInfo, err error) {
	b, err := dbServer.CommandPost("jobs", &t)

	if err != nil {
		return
	}

	err = json.NewDecoder(b).Decode(&job)
	return
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

	job.Resource = prio.Sort()[0]

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
