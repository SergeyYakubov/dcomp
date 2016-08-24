// Package provides an infrastructure for communications with a generic server.
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Server struct {
	Host string
	Port int
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
		s = "/" + s + "/"
	}
	return fmt.Sprintf("http://%s:%d%s", srv.Host, srv.Port, s)
}

// CommandPost issues the POST command to srv. data should be JSON-encodable. Returns response body or error
func (srv *Server) CommandPost(path string, data interface{}) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return nil, err
	}
	res, err := http.Post(srv.url(path), "application/json; charset=utf-8", b)

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

// CommandGet issues the GET command to srv. Returns response body or error
func (srv *Server) CommandGet(path string) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	res, err := http.Get(srv.url(path))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	io.Copy(b, res.Body)

	if res.StatusCode != http.StatusOK {
		err = errors.New(b.String())
		return nil, err
	}

	return b, nil
}

// CommandPost issues the DELETE command to srv. Returns response body or error
func (srv *Server) CommandDelete(path string) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	req, err := http.NewRequest(http.MethodDelete, srv.url(path), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	io.Copy(b, res.Body)

	if res.StatusCode != http.StatusOK {
		err = errors.New(b.String())
		return nil, err
	}

	return b, nil
}

// CommandPatch issues the PATCH command to srv. Returns response body or error
func (srv *Server) CommandPatch(path string, data interface{}) (err error) {
	b := new(bytes.Buffer)

	if err := json.NewEncoder(b).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, srv.url(path), b)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		io.Copy(b, res.Body)
		err = errors.New(b.String())
		return err
	}

	return nil
}
