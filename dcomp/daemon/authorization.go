package daemon

import (
	"net/http"

	"encoding/json"

	"github.com/gorilla/context"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

func ProcessUserAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, _, err := server.ExtractAuthInfo(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		resp, err := authorize(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		context.Set(r, "authorizationResponce", &resp)
		defer context.Clear(r)
		fn(w, r)
	}
}

func authorize(r *http.Request) (server.AuthorizationResponce, error) {

	var req server.AuthorizationRequest
	var resp server.AuthorizationResponce

	req.Command = r.Method
	req.Token = r.Header.Get("Authorization")
	req.URL = r.URL.RawPath

	b, err := authServer.CommandPost("authorize"+"/", &req)

	if err != nil {
		return resp, err
	}

	if err := json.NewDecoder(b).Decode(&resp); err != nil {
		return resp, err
	}

	return resp, nil
}
