package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []structs.JobInfo
	if err := database.GetAllRecords(&jobs); err != nil {
		http.Error(w, "cannot retrieve database job info", http.StatusBadRequest)
		return
	}
	if len(jobs) == 0 {
		fmt.Fprint(w, "No jobs found")
		return
	}

	for _, v := range jobs {
		fmt.Fprintf(w, "Job: %s Image: %s NCPUs %d\n", v.Id, v.ImageName, v.NCPUs)
	}
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	var jobs []structs.JobInfo

	if err := database.GetRecordById(jobID, &jobs); err != nil {
		http.Error(w, "cannot retrieve database job info", http.StatusBadRequest)
		return
	}
	if len(jobs) == 0 {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Job: %s Image: %s NCPUs %d\n", jobs[0].Id, jobs[0].ImageName, jobs[0].NCPUs)

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
