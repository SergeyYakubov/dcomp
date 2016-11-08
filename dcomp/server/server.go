// Package provides an infrastructure for communications with a generic server.
package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"github.com/dcomp/dcomp/utils"
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
	Host string
	Port int
	Key  string
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
	return fmt.Sprintf("http://%s:%d%s", srv.Host, srv.Port, s)
}

func (srv *Server) addAuthorizationHeader(r *http.Request, s string) {
	if srv.Key == "" {
		return
	}

	u := utils.StripURL(r.URL)
	mac := hmac.New(sha256.New, []byte(srv.Key))
	mac.Write([]byte(u))
	if s != "" {
		mac.Write([]byte(s))
	}
	sha := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	r.Header.Add("Authorization", sha)
}

// CommandDelete issues the http command to srv. Returns response body or error
func (srv *Server) httpCommand(method string, path string, data interface{}) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)
	if data != nil {
		if err := json.NewEncoder(b).Encode(data); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, srv.url(path), b)
	if err != nil {
		return nil, err
	}

	req.Close = true

	srv.addAuthorizationHeader(req, "")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	io.Copy(b, res.Body)

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		err = errors.New(b.String())
		return nil, err
	}

	return b, nil
}

// CommandPost issues the POST command to srv. data should be JSON-encodable. Returns response body or error
func (srv *Server) CommandPost(path string, data interface{}) (b *bytes.Buffer, err error) {
	return srv.httpCommand(http.MethodPost, path, data)
}

// CommandGet issues the GET command to srv. Returns response body or error
func (srv *Server) CommandGet(path string) (b *bytes.Buffer, err error) {
	return srv.httpCommand(http.MethodGet, path, nil)
}

// CommandDelete issues the DELETE command to srv. Returns response body or error
func (srv *Server) CommandDelete(path string) (b *bytes.Buffer, err error) {
	return srv.httpCommand(http.MethodDelete, path, nil)
}

// CommandPatch issues the PATCH command to srv. Returns response body or error
func (srv *Server) CommandPatch(path string, data interface{}) (err error) {
	_, err = srv.httpCommand(http.MethodPatch, path, data)
	return
}
