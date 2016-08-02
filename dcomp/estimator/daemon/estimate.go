package daemon

import (
	"net/http"

	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func EstimateJob(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	var t structs.JobDescription

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	//	b := new(bytes.Buffer)
	//	json.NewEncoder(b).Encode(t)
	//	w.Write(b.Bytes())
}
