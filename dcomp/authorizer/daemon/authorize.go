package daemon

import (
	"net/http"
	"time"

	"bytes"
	"encoding/json"

	"crypto/tls"
	"fmt"
	"strings"

	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"gopkg.in/ldap.v2"
)

func routeAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-type", "application/json")

	if r.Body == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var t server.AuthorizationRequest

	d := json.NewDecoder(r.Body)
	if d.Decode(&t) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	resp, err := authorize(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(resp)
	w.Write(b.Bytes())
}

// authorize checks user authorization and returns responce
func authorize(req server.AuthorizationRequest) (server.AuthorizationResponce, error) {

	var resp server.AuthorizationResponce

	atype, atoken, err := server.SplitAuthToken(req.Token)

	if err != nil {
		resp.Status = http.StatusUnauthorized
		resp.StatusText = "wrong token"
		return resp, nil
	}

	if !c.authorizationAllowed(atype) {
		resp.Status = http.StatusUnauthorized
		resp.StatusText = "wrong auth type"
		return resp, nil
	}

	switch atype {
	case "None":
		return AuthorizeWithToken(atoken)
	case "Negotiate":
		if gssAPIContext == nil {
			err = errors.New("gssAPIContext not defined")
			return resp, err
		}
		resp.UserName, err = gssAPIContext.ParseToken(atoken)
		if err != nil {
			resp.StatusText = "Wrong username or password"
			resp.Status = http.StatusUnauthorized
			return resp, nil
		}
		resp.Status = http.StatusOK
		resp.Token = atoken
		return resp, nil
	case "Bearer":
		claim, err := extractJWTTokenClaim(atoken)
		if err != nil {
			resp.StatusText = "token invalid"
			resp.Status = http.StatusUnauthorized
			return resp, nil
		}

		return AuthorizeWithToken(claim.User)

	case "Basic":
		user, errl := ldap_login(atoken)
		if errl != nil {
			resp.StatusText = "Wrong username or password"
			resp.Status = http.StatusUnauthorized
			return resp, nil
		}
		return AuthorizeWithToken(user)
	}
	resp.StatusText = "wrong auth type"
	resp.Status = http.StatusUnauthorized
	return resp, nil
}

func ldap_login(token string) (string, error) {

	uEnc, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return "", err
	}
	creds := strings.Split(string(uEnc), ":")

	if len(creds) != 2 {
		return "", errors.New("wrong token")
	}

	username := creds[0]
	password := creds[1]

	l, err := ldap.Dial("tcp", c.Ldap.Host)
	if err != nil {
		return "", err
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		c.Ldap.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", errors.New("User does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	fullname := sr.Entries[0].GetAttributeValue("displayName")
	if fullname != "" {
		username = fullname
	}

	// Bind as the user to verify their password
	return username, l.Bind(userdn, password)

}

type tokenClaims struct {
	User string
}

func AuthorizeWithToken(user string) (resp server.AuthorizationResponce, err error) {
	resp.Token, err = createJWTToken(user)
	if err != nil {
		return
	}

	resp.UserName = user
	resp.Status = http.StatusOK
	resp.ValidityTime = c.Tokenduration
	return

}

func createJWTToken(user string) (string, error) {

	var claims server.CustomClaims
	var extraClaim tokenClaims

	extraClaim.User = user
	claims.ExtraClaims = &extraClaim
	claims.Duration = time.Duration(c.Tokenduration) * time.Minute

	token := server.NewJWTAuth(c.Daemon.Key)

	return token.GenerateToken(&claims)
}

func extractJWTTokenClaim(token string) (tokenClaims, error) {

	claims, ok := server.CheckJWTToken(token, c.Daemon.Key)
	if !ok {
		return tokenClaims{}, errors.New("token invalid")
	}
	customClaim, ok := claims.(*server.CustomClaims)
	if !ok {
		return tokenClaims{}, errors.New("token invalid")
	}
	var tokenClaim tokenClaims
	err := utils.MapToStruct(customClaim.ExtraClaims.(map[string]interface{}), &tokenClaim)

	if err != nil {
		return tokenClaims{}, errors.New("token invalid")
	}

	return tokenClaim, nil
}
