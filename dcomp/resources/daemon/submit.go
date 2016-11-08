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

	err := resource.SubmitJob(t)
	if err != nil {
		http.Error(w, "cannot submit job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//	b := new(bytes.Buffer)
	//	json.NewEncoder(b).Encode(res)
	//	w.Write(b.Bytes())

}
