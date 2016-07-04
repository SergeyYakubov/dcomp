package daemon

import (
	"encoding/json"
	"io"
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/common_structs"
)

func decodeJob(r io.Reader) (commonStructs.JobDescription, bool) {

	var t commonStructs.JobDescription

	if r == nil {
		return t, false
	}

	decoder := json.NewDecoder(r)

	if decoder.Decode(&t) != nil || t.Check() != nil {
		return t, false
	}

	return t, true
}

func SubmitJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	t, ok := decodeJob(r.Body)

	if !ok {
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
