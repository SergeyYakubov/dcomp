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

		b, status, err := authorize(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if status != http.StatusOK {
			w.WriteHeader(status)
			w.Write(b.Bytes())
			return
		}

		var resp server.AuthorizationResponce
		if err := json.NewDecoder(b).Decode(&resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "authorizationResponce", &resp)
		fn(w, r.WithContext(ctx))
	}
}

func authorize(r *http.Request) (*bytes.Buffer, int, error) {

	var req server.AuthorizationRequest

	req.Command = r.Method
	req.Token = r.Header.Get("Authorization")
	req.URL = r.URL.RawPath

	return authServer.CommandPost("authorize"+"/", &req)
}
