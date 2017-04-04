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

		/*		if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers",
						"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				}
				// Stop here if its Preflighted OPTIONS request
				if r.Method == "OPTIONS" {
					return
				}

		*/

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

		if resp.Status != http.StatusOK {
			http.Error(w, resp.StatusText, resp.Status)
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
