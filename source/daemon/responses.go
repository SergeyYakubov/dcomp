package daemon

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllJobs(w http.ResponseWriter, r *http.Request) {

}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	fmt.Fprintln(w, "jobID show:", jobID)
}
