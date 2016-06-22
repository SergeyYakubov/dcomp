package daemon

import (
	"../common_structs"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Fprintf(w, "%q\n", err)
		fmt.Fprintf(w, "%v\n", r)

	}
}
