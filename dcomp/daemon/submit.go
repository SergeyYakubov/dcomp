package daemon

import (
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func SubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobDescription

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b, err := DBServer.CommandPost("jobs", &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		go findResourceAndSubmitJob(t, b.String())
		w.WriteHeader(http.StatusCreated)
		w.Write(b.Bytes())
	}
}

func findResourceAndSubmitJob(t structs.JobDescription, id string) {

}
