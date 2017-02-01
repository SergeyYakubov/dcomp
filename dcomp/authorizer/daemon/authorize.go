package daemon

import (
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

func unAuthorizedResponce(msg string, w http.ResponseWriter) {

	var resp server.AuthorizationResponce

	resp.Authorization = make([]string, len(c.Authorization))
	copy(resp.Authorization, c.Authorization)
	resp.ErrorMsg = msg

	w.WriteHeader(http.StatusUnauthorized)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&resp)
	w.Write(b.Bytes())
	return
}

func routeAuthorizeRequest(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("Content-type", "application/json")

	if r.Body == nil {
		unAuthorizedResponce("bad request", w)
		return

	}

	var t server.AuthorizationRequest

	d := json.NewDecoder(r.Body)
	if d.Decode(&t) != nil {
		unAuthorizedResponce("bad request", w)
		return

	}

	resp, err := authorize(t)
	if err != nil {
		unAuthorizedResponce(err.Error(), w)
		return
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	w.Write(b.Bytes())
}

// authorize checks user authorization and returns responce
func authorize(req server.AuthorizationRequest) (resp server.AuthorizationResponce, err error) {
	atype, atoken, err := server.SplitAuthToken(req.Token)
	if err != nil {
		return
	}

	if !c.authorizationAllowed(atype) {
		err = errors.New("wrong auth type")
		return
	}

	switch atype {
	case "Negotiate":
		if gssAPIContext == nil {
			err = errors.New("gssAPIContext not defined")
			return
		}
		resp.UserName, err = gssAPIContext.ParseToken(atoken)
	case "Basic":
		resp.UserName = atoken
	}
	return
}
