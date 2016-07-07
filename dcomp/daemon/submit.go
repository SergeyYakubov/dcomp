package daemon

import (
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/common_structs"
)

func SubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t commonStructs.JobDescription

	if ok := commonStructs.DecodeStruct(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b, err := DBServer.PostCommand("jobs", &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		go findResourceAndSubmitJob(t, b.String())
		w.WriteHeader(http.StatusCreated)
		w.Write(b.Bytes())
	}
}

func findResourceAndSubmitJob(t commonStructs.JobDescription, id string) {

}
