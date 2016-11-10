// Package provides an infrastructure for communications with a generic server.
package server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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
	auth Auth
	Tls  bool
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

	token, err := srv.auth.GenerateToken(r)
	if err != nil {
		log.Print("cannot generat auth token " + err.Error())
		return
	}

	r.Header.Add("Authorization", token)
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

	srv.addAuthorizationHeader(req)

	var client *http.Client
	if srv.Tls {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}

	res, err := client.Do(req)

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
