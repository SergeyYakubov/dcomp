package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"stash.desy.de/scm/dc/main.git/dcomp/common_structs"
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

	var t commonStructs.JobInfo
	if decoder.Decode(&t) == nil && t.Check() == nil {
		t.Id = "1"
		t.Status = 1
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(t)
		w.WriteHeader(http.StatusCreated)
		w.Write(b.Bytes())
	} else {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
