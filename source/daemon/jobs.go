package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"stash.desy.de/scm/dc/common_structs"
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
	decoder := json.NewDecoder(r.Body)

	var t commonStructs.JobDescription
	if decoder.Decode(&t) == nil && t.Check() == nil {
		fmt.Fprintf(w, "Job submitted\n")
	} else {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
