package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func sendJobs(w http.ResponseWriter, jobs []structs.JobInfo, allowempty bool) {
	if len(jobs) == 0 && allowempty {
		return
	}

	if len(jobs) == 0 {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(jobs); err != nil {
		http.Error(w, "cannot retrieve database job info", http.StatusInternalServerError)
		return
	}
	w.Write(b.Bytes())

}

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []structs.JobInfo
	if err := database.GetAllRecords(&jobs); err != nil {
		http.Error(w, "cannot retrieve database job info", http.StatusBadRequest)
		return
	}

	sendJobs(w, jobs, true)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	var jobs []structs.JobInfo

	if err := database.GetRecordById(jobID, &jobs); err != nil {
		http.Error(w, "cannot retrieve database job info", http.StatusBadRequest)
		return
	}

	sendJobs(w, jobs, false)

}

func SubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobInfo

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	t.Status = 1
	id, err := database.CreateRecord(t)
	if err != nil {
		http.Error(w, "bad request "+err.Error(), http.StatusBadRequest)
		return
	}
	t.Id = id
	w.WriteHeader(http.StatusCreated)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(t)
	w.Write(b.Bytes())
}
