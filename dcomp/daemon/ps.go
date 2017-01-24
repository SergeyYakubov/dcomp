package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"net/url"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
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
		if jobs[i].Status != structs.StatusFinished && jobs[i].Status != structs.StatusWaitData {
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

func getJobLog(job structs.JobInfo, u *url.URL) (*bytes.Buffer, error) {
	res := resources[job.Resource]
	cmd := u.Path + "?" + u.RawQuery
	return res.Server.CommandGet(cmd)
}

func routeGetJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jobLog := r.URL.Query().Get("log")

	if jobLog == "true" {
		b, err := getJobLog(job, r.URL)
		if err != nil {
			http.Error(w, "cannot retrieve job log: "+err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b.Bytes())
		return
	}

	jobs := []structs.JobInfo{job}
	updateJobs(jobs)
	sendJobs(w, jobs, false)
}
