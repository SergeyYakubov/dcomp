package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/gorilla/mux"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func sendJobs(w http.ResponseWriter, jobs []structs.JobInfo, allowempty bool) {
	if len(jobs) == 0 && allowempty {
		w.WriteHeader(http.StatusOK)
		return
	}

	if len(jobs) == 0 {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(jobs); err != nil {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())

}

func updateJobs(jobs []structs.JobInfo) {
	for i, _ := range jobs {
		if jobs[i].Status != structs.StatusFinished {
			updateJobsStatusFromResources(&jobs[i])
		}
	}
}

func updateJobsStatusFromResources(job *structs.JobInfo) {

	res := resources[job.Resource]
	// update database
	defer db.PatchRecord(job.Id, job)

	b, err := res.Server.CommandGet("jobs" + "/" + job.Id)

	if err != nil {
		job.Status = structs.StatusErrorFromResource
		job.Message = err.Error()
		return
	}

	var status structs.JobStatus
	if err := json.NewDecoder(b).Decode(&status); err != nil {
		job.Status = structs.StatusErrorFromResource
		job.Message = err.Error()
		return
	}

	job.JobStatus = status

	return
}

func pickNotFinished(jobs []structs.JobInfo) (res []structs.JobInfo) {
	res = make([]structs.JobInfo, 0)
	for _, job := range jobs {
		if job.Status != structs.StatusFinished {
			res = append(res, job)
		}
	}
	return
}

func routeGetAllJobs(w http.ResponseWriter, r *http.Request) {

	var jobs []structs.JobInfo

	if err := db.GetAllRecords(&jobs); err != nil {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusBadRequest)
		return
	}

	showFinished := r.URL.Query().Get("finished")
	if showFinished != "true" {
		jobs = pickNotFinished(jobs)
	}
	updateJobs(jobs)
	sendJobs(w, jobs, true)
}

func routeGetJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	var jobs []structs.JobInfo
	if err := db.GetRecordByID(jobID, &jobs); err != nil {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusNotFound)
		return
	}
	updateJobs(jobs)
	sendJobs(w, jobs, false)
}
