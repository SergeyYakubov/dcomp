package daemon

import (
	"net/http"

	"encoding/json"

	"errors"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func writeSubmitResponce(w http.ResponseWriter, r *http.Request, job structs.JobInfo) {
	if job.Status == structs.StatusSubmitted {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&job)
		return
	}
	if job.Status == structs.StatusWaitData {

		err := writeJWTToken(w, r, job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}

	http.Error(w, "unknow job status", http.StatusInternalServerError)
	return
}

func routeSubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobDescription

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	job, err := trySubmitJob(user, t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSubmitResponce(w, r, job)

	return
}

func routeReleaseJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = submitToResource(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	job.Status = structs.StatusSubmitted
	if err := modifyJobInDatabase(job.Id, &job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSubmitResponce(w, r, job)
	return

}

func trySubmitJob(user string, t structs.JobDescription) (job structs.JobInfo, err error) {

	job.JobUser = user
	job, err = addJobToDatabase(t)
	if err != nil {
		return
	}

	prio, err := findResources(t)
	if err != nil {
		return
	}

	job.Resource, err = checkResources(job, prio.Sort())
	if err != nil {
		return
	}

	if t.NeedData() {
		job.Status = structs.StatusWaitData
	} else {
		if err = submitToResource(job); err != nil {
			return
		}
		job.Status = structs.StatusSubmitted
	}

	err = modifyJobInDatabase(job.Id, &job)
	return
}

func modifyJobInDatabase(id string, data interface{}) error {
	return db.PatchRecord(id, data)
}

func addJobToDatabase(t structs.JobDescription) (job structs.JobInfo, err error) {

	job.JobDescription = t

	job.Id, err = db.CreateRecord("", &job)
	if err != nil {
		return
	}

	return
}

func checkResources(job structs.JobInfo, prio []string) (res string, err error) {
	for i := range prio {
		r, ok := resources[prio[i]]
		if ok {
			_, e := r.Server.CommandPost("jobs/?checkonly=true", job)
			if e == nil {
				res = prio[i]
				return
			}
		}
	}

	err = errors.New("no resource available")
	return
}

func submitToResource(job structs.JobInfo) (err error) {
	r, ok := resources[job.Resource]
	if !ok {
		err = errors.New("Resource unvailable " + job.Resource)
	}
	_, err = r.Server.CommandPost("jobs", job)
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
