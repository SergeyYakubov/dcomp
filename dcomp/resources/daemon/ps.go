package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func writeStatus(w http.ResponseWriter, status structs.JobStatus) {
	w.WriteHeader(http.StatusOK)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(status)
	w.Write(b.Bytes())
}

func writeLogs(w http.ResponseWriter, id string, compress bool) {

	b, err := resource.GetLogs(id, compress)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())

}

func routeGetJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	jobLog := r.URL.Query().Get("log")

	if jobLog == "true" {
		logCompress := r.URL.Query().Get("compress")
		writeLogs(w, jobID, logCompress == "true")
	} else {
		status, err := resource.GetJobStatus(jobID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		writeStatus(w, status)
	}

}
