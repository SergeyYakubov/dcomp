package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/gorilla/mux"
)

func routeGetJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	status, err := resource.GetJob(jobID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(status)
	w.Write(b.Bytes())

}
