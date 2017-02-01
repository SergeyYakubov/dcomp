package daemon

import (
	"net/http"

	"encoding/json"

	"bytes"
	"context"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

func ProcessUserAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//		_, _, err := server.ExtractAuthInfo(r)

		//		if err != nil {
		//			http.Error(w, "authorization failed: "+err.Error(), http.StatusUnauthorized)
		//			return
		//		}

		resp, status, err := authorize(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if status == http.StatusUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(&resp)
			w.Write(b.Bytes())
			return
		}
		ctx := context.WithValue(r.Context(), "authorizationResponce", &resp)
		fn(w, r.WithContext(ctx))
	}
}

func authorize(r *http.Request) (server.AuthorizationResponce, int, error) {

	var req server.AuthorizationRequest
	var resp server.AuthorizationResponce

	req.Command = r.Method
	req.Token = r.Header.Get("Authorization")
	req.URL = r.URL.RawPath

	b, status, err := authServer.CommandPost("authorize"+"/", &req)
	if err != nil {
		return resp, -1, err
	}

	if err := json.NewDecoder(b).Decode(&resp); err != nil {
		return resp, -1, err
	}

	return resp, status, err
}
