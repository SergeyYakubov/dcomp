package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

func routeAuthorizeRequest(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	if r.Body == nil {
		http.Error(w, "bad request", http.StatusUnauthorized)
		return

	}

	var t server.AuthorizationRequest

	d := json.NewDecoder(r.Body)
	if d.Decode(&t) != nil {
		http.Error(w, "bad request", http.StatusUnauthorized)
		return

	}

	resp, err := authorize(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	w.Write(b.Bytes())
}

// authorize checks user authorization and returns responce
func authorize(req server.AuthorizationRequest) (resp server.AuthorizationResponce, err error) {
	at, au, err := server.SplitAuthToken(req.Token)
	if err != nil {
		return
	}

	if at != "Basic" {
		err = errors.New("wrong auth type")
		return
	}
	resp.UserName = au
	return
}
