// Package provides an infrastructure for communications with a generic server.
package server

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type Server struct {
	Host            string
	Port            int
	auth            Auth
	Tls             bool
	alternativeAuth []Auth
}

func (srv *Server) AddAlternativeAuth(a Auth) {
	if srv.alternativeAuth == nil {
		srv.alternativeAuth = make([]Auth, 0)
	}
	for _, t := range srv.alternativeAuth {
		if t.Name() == a.Name() {
			return
		}
	}
	srv.alternativeAuth = append(srv.alternativeAuth, a)
}

func (srv *Server) SetAuth(a Auth) {
	srv.auth = a
}

func (srv *Server) GetAuth() Auth {
	return srv.auth
}

// ParseUrl extacts host anf port from a string and sets corresponding structure fields
func (srv *Server) parseUrl(s string) {
	u, _ := url.Parse(s)
	host, port, _ := net.SplitHostPort(u.Host)
	srv.Host = host
	srv.Port, _ = strconv.Atoi(port)
}

//
func (srv *Server) FullName() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func (srv *Server) url(s string) string {
	if s != "" {
		s = strings.TrimSpace(s)
		s = strings.TrimLeft(s, "/")
		s = strings.TrimRight(s, "/")
		s = "/" + s
		if !strings.ContainsRune(s, '?') {
			s += "/"
		}
	}
	protocol := "http"
	if srv.Tls {
		protocol += "s"
	}
	return fmt.Sprintf(protocol+"://%s:%d%s", srv.Host, srv.Port, s)
}

func (srv *Server) addAuthorizationHeader(r *http.Request) {
	if srv.auth == nil {
		return
	}

	claims := CustomClaims{ExtraClaims: r}
	token, err := srv.auth.GenerateToken(&claims)
	if err != nil {
		log.Print("cannot generate auth token: " + err.Error())
		return
	}

	r.Header.Add("Authorization", token)
}

func (srv *Server) newClient() (client *http.Client) {
	if srv.Tls {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}

	return
}

func (srv *Server) UploadData(urlpath string, destname string, data io.Reader,
	size int64, mode os.FileMode) (b *bytes.Buffer, err error) {

	req, err := http.NewRequest("POST", srv.url(urlpath), data)
	if err != nil {
		return nil, err
	}

	srv.addAuthorizationHeader(req)

	cd := "attachment; filename=" + url.QueryEscape(destname)

	req.Header.Set("Content-Disposition", cd)
	req.Header.Set("Content-Type", "application/octet-stream")

	m := new(bytes.Buffer)
	binary.Write(m, binary.LittleEndian, mode)
	req.Header.Set("X-Content-Mode", url.QueryEscape(m.String()))

	client := srv.newClient()

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	b = new(bytes.Buffer)

	defer res.Body.Close()

	io.Copy(b, res.Body)

	if res.StatusCode != http.StatusCreated {
		err = errors.New(b.String())
		return nil, err
	}
	return b, nil
}

func (srv *Server) httpCommandWithAuth(method string, path string, data interface{}) (*bytes.Buffer,
	int, error) {

	b := new(bytes.Buffer)
	if data != nil {
		if err := json.NewEncoder(b).Encode(data); err != nil {
			return nil, -1, err
		}
	}

	req, err := http.NewRequest(method, srv.url(path), b)
	if err != nil {
		return nil, -1, err
	}

	req.Close = true
	srv.addAuthorizationHeader(req)

	client := srv.newClient()
	res, err := client.Do(req)

	if err != nil {
		return nil, -1, err
	}

	defer res.Body.Close()
	io.Copy(b, res.Body)

	return b, res.StatusCode, nil
}

func (srv *Server) authTokensToTry() (tokens []Auth) {
	tokens = make([]Auth, 0)
	tokens = append(tokens, srv.auth)
	f := func(token Auth) {
		for _, t := range tokens {
			if t == nil {
				continue
			}
			if t.Name() == token.Name() {
				return
			}
		}
		tokens = append(tokens, token)
	}

	for _, t := range srv.alternativeAuth {
		f(t)
	}

	return
}

func (srv *Server) httpCommand(method string, path string, data interface{}) (b *bytes.Buffer, status int, err error) {

	tryAuthTokens := srv.authTokensToTry()
	iniToken := srv.auth
	defer func() { srv.auth = iniToken }()

	for _, authToken := range tryAuthTokens {
		srv.auth = authToken

		b, status, err = srv.httpCommandWithAuth(method, path, data)
		if err != nil || status != http.StatusUnauthorized {
			return
		}
	}

	var tokens string
	for _, t := range tryAuthTokens {
		if t != nil {
			tokens += " " + t.Name()
		} else {
			tokens += " None"
		}
	}
	b = new(bytes.Buffer)
	b.WriteString("Cannot authorize with methods:" + tokens)
	return
}

// CommandPost issues the POST command to srv. data should be JSON-encodable. Returns response body or error
func (srv *Server) CommandPost(path string, data interface{}) (b *bytes.Buffer, status int, err error) {
	return srv.httpCommand(http.MethodPost, path, data)
}

// CommandGet issues the GET command to srv. Returns response body or error
func (srv *Server) CommandGet(path string) (b *bytes.Buffer, status int, err error) {
	return srv.httpCommand(http.MethodGet, path, nil)
}

// CommandDelete issues the DELETE command to srv. Returns response body or error
func (srv *Server) CommandDelete(path string) (b *bytes.Buffer, status int, err error) {
	return srv.httpCommand(http.MethodDelete, path, nil)
}

// CommandPatch issues the PATCH command to srv. Returns response body or error
func (srv *Server) CommandPatch(path string, data interface{}) (*bytes.Buffer, int, error) {
	return srv.httpCommand(http.MethodPatch, path, data)
}
