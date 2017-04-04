package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

func routeLogin(w http.ResponseWriter, r *http.Request) {

	resp := r.Context().Value("authorizationResponce").(*server.AuthorizationResponce)
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(resp); err != nil {
		http.Error(w, "cannot encode authorization responce: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())

}
