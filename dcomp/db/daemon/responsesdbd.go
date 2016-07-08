package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"stash.desy.de/scm/dc/main.git/dcomp/db/database"
)

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	fmt.Fprintln(w, "jobID show:", jobID)
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
