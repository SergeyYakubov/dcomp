package daemon

import (
	"net/http"
	"time"

	"bytes"
	"encoding/json"

	"net/url"

	"fmt"
	"strings"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"strconv"
)

func sendJobs(w http.ResponseWriter, jobs []structs.JobInfo, allowempty bool) {
	if len(jobs) == 0 && !allowempty {
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
		if jobs[i].Status != structs.StatusFinished &&
			jobs[i].Status != structs.StatusWaitData &&
			jobs[i].Status != structs.StatusSubmissionFailed &&
			jobs[i].Status != structs.StatusUserDataCopied &&
			jobs[i].Status != structs.StatusDataCopyFailed {

			updateJobsStatusFromResources(&jobs[i])
		}
		jobs[i].UpdateStatusString()
	}
}

func updateJobsStatusFromResources(job *structs.JobInfo) {

	res := resources[job.Resource]
	// update database
	defer db.PatchRecord(job.Id, job)

	b, httpstatus, err := res.Server.CommandGet("jobs" + "/" + job.Id)

	if httpstatus == http.StatusNotFound {
		return
	}

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

func filterJobs(jobs []structs.JobInfo, filter url.Values) (res []structs.JobInfo) {
	res = make([]structs.JobInfo, 0)
	for _, job := range jobs {

		if filter.Get("finishedOnly") == "true" && job.Status != structs.StatusFinished {
			continue
		}

		if filter.Get("notFinishedOnly") == "true" && job.Status == structs.StatusFinished {
			continue
		}

		if filter.Get("keyword") != "" {
			str := fmt.Sprintf("%v", job)
			if !strings.Contains(str, filter.Get("keyword")) {
				continue
			}

		}

		from := filter.Get("from")
		to := filter.Get("to")
		ftmstring := "2006-01-02"
		if from != "" && to != "" {
			timeFrom, err1 := time.Parse(ftmstring, from)
			timeTo, err2 := time.Parse(ftmstring, to)
			jobTime, err3 := utils.StringToTime(job.SubmitTime)
			if err1 == nil && err2 == nil && err3 == nil {
				if jobTime.Before(timeFrom) || jobTime.After(timeTo) {
					continue
				}
			}
		} else if filter.Get("last") != "" {
			days, err := strconv.Atoi(filter.Get("last"))
			jobTime, err1 := utils.StringToTime(job.SubmitTime)
			if err == nil && err1 == nil {
				timelast := time.Now().Add(-time.Duration(days*24) * time.Hour)

				if jobTime.Before(timelast) {
					continue
				}

			}

		}

		res = append(res, job)
	}
	return
}

func routeGetAllJobs(w http.ResponseWriter, r *http.Request) {

	var jobs []structs.JobInfo

	if err := db.GetAllRecords(&jobs); err != nil {
		http.Error(w, "cannot retrieve database job info: "+err.Error(), http.StatusBadRequest)
		return
	}

	updateJobs(jobs)

	filter := r.URL.Query()
	jobs = filterJobs(jobs, filter)

	sendJobs(w, jobs, true)
}

func getJobLog(job structs.JobInfo, u *url.URL) (b *bytes.Buffer, err error) {
	res := resources[job.Resource]
	cmd := u.Path + "?" + u.RawQuery
	b, _, err = res.Server.CommandGet(cmd)
	return
}

func routeGetJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromRequest(r)

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
