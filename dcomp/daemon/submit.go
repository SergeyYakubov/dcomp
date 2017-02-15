package daemon

import (
	"net/http"
	"time"

	"encoding/json"

	"errors"

	"strings"

	"fmt"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func writeSubmitResponce(w http.ResponseWriter, job structs.JobInfo) {
	if job.Status == structs.StatusSubmitted {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&job)
		return
	}

	if job.Status == structs.StatusUserDataCopied {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(&job)
		return
	}

	if job.Status == structs.StatusWaitData {
		w.WriteHeader(http.StatusAccepted)
		if job.NeedUserDataUpload() {
			err := writeJWTToken(w, job)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			json.NewEncoder(w).Encode(&job)
		}
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

	writeSubmitResponce(w, job)

	return
}

func routeReleaseJob(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// job will be submitted automatically when internal data is copied
	// we need to set flag in database for this and exit
	if job.NeedInternalDataCopy() {
		if err := setJobStatus(&job, structs.StatusUserDataCopied); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeSubmitResponce(w, job)
		return
	}

	if err = submitToResource(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := setJobStatus(&job, structs.StatusSubmitted); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeSubmitResponce(w, job)
	return
}

func failedDataCopy(job structs.JobInfo, inierr error) error {
	if errdb := setJobStatus(&job, structs.StatusDataCopyFailed); errdb != nil {
		return errors.New(inierr.Error() + errdb.Error())
	}
	return inierr

}

func submitSingleCopyDataRequest(job structs.JobInfo, fi structs.FileCopyInfo, errchan chan error) {
	sourceJobID := fi.Source
	sourceJob, err := GetJobFromDatabase(sourceJobID)
	if err != nil {
		errchan <- failedDataCopy(job, err)
		return
	}
	if sourceJob.Resource != job.Resource {
		errchan <- failedDataCopy(job, errors.New("Cannot copy/mount from another resource"))
		return
	}
	token, err := createJWT(job)
	srv := resources[job.Resource].DataManager
	srv.SetAuth(server.NewExternalAuth(token))
	b, status, err := srv.CommandPost("jobfile/"+job.Id+"/?mode=mount", &fi)

	if err != nil {
		errchan <- err
		return
	}

	if status != http.StatusCreated {
		errchan <- errors.New(b.String())
		return
	}

	errchan <- nil

	return

}

func waitJobStatus(jobID string, status int, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		job, err := GetJobFromDatabase(jobID)
		if err != nil {
			return err
		}
		if job.Status == status {
			return nil
		}

		time.Sleep(time.Second)
	}
	return errors.New("Timeout waiting status " + structs.JobStatusExplained[status])

}

func waitRequests(n int, errchan chan error) (err error) {
	// wait for single requests to finish, accumulate error
	for i := 0; i < n; i++ {
		err_current := <-errchan
		if err_current != nil {
			if err == nil {
				err = err_current
			} else {
				err = errors.New(err.Error() + err_current.Error())
			}
		}
	}
	return err
}

func submitAfterCopyDataRequest(job structs.JobInfo, waitUserData bool) error {

	errchan := make(chan error)

	nrequests := 0

	for _, fi := range job.FilesToMount {
		go submitSingleCopyDataRequest(job, fi, errchan)
		nrequests++
	}

	if err := waitRequests(nrequests, errchan); err != nil {
		setJobStatus(&job, structs.StatusDataCopyFailed, err.Error())
		return err
	}

	if waitUserData {

		timeout := time.Hour * 12
		if err := waitJobStatus(job.Id, structs.StatusUserDataCopied, timeout); err != nil {
			setJobStatus(&job, structs.StatusDataCopyFailed, err.Error())
			fmt.Println(err)
			return err
		}

	}

	if err := submitToResource(job); err != nil {
		setJobStatus(&job, structs.StatusSubmissionFailed, err.Error())
		return err
	}

	return setJobStatus(&job, structs.StatusSubmitted)
}

func trySubmitJob(user string, t structs.JobDescription) (job structs.JobInfo, err error) {

	job, err = addJobToDatabase(t)
	job.JobUser = user

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
	if err = modifyJobInDatabase(job.Id, &job); err != nil {
		return
	}

	if t.NeedInternalDataCopy() || t.NeedUserDataUpload() {
		if err = setJobStatus(&job, structs.StatusWaitData); err != nil {
			return
		}
		if t.NeedInternalDataCopy() {
			go submitAfterCopyDataRequest(job, t.NeedUserDataUpload())
		}
		return
	}

	if err = submitToResource(job); err != nil {
		return
	}

	err = setJobStatus(&job, structs.StatusSubmitted)
	return
}

func setJobStatus(job *structs.JobInfo, status int, message ...string) error {

	job.Status = status
	if len(message) > 0 {
		job.Message = message[0]
	}
	data := struct {
		JobStatus structs.JobStatus
	}{structs.JobStatus{Status: status, Message: job.Message}}

	return db.PatchRecord(job.Id, &data)
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
			_, _, e := r.Server.CommandPost("jobs/?checkonly=true", job)
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
	_, _, err = r.Server.CommandPost("jobs", &job)
	return
}

func findResources(t structs.JobDescription) (prio structs.ResourcePrio, err error) {
	if t.Resource != "" {
		prio = make(structs.ResourcePrio)
		prio[strings.ToLower(t.Resource)] = 100
		return
	}
	b, _, err := estimatorServer.CommandPost("estimations", &t)
	if err != nil {
		return
	}
	err = json.NewDecoder(b).Decode(&prio)
	return
}
