package daemon

import (
	"net/http"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func routeSubmitJob(w http.ResponseWriter, r *http.Request) {

	var t structs.JobInfo

	if ok := structs.Decode(r.Body, &t); !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	q := r.URL.Query().Get("checkonly")

	success := http.StatusCreated
	checkonly := false
	if q == "true" {
		success = http.StatusOK
		checkonly = true

	}

	err := resource.SubmitJob(t, checkonly)
	if err != nil {
		http.Error(w, "cannot submit job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(success)
	//	b := new(bytes.Buffer)
	//	json.NewEncoder(b).Encode(res)
	//	w.Write(b.Bytes())

}
