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

func (srv *Server) ParseUrl(s string) {
	u, _ := url.Parse(s)
	host, port, _ := net.SplitHostPort(u.Host)
	srv.Host = host
	srv.Port, _ = strconv.Atoi(port)

}

func (srv *Server) HostPort() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func (srv *Server) Url(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, "/")
	s = strings.TrimRight(s, "/")
	s = "/" + s + "/"
	return fmt.Sprintf("http://%s:%d%s", srv.Host, srv.Port, s)
}

func (srv *Server) PostCommand(path string, data interface{}) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return nil, err
	}
	res, err := http.Post(srv.Url(path), "application/json; charset=utf-8", b)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	io.Copy(b, res.Body)

	if res.StatusCode != http.StatusCreated {
		err = errors.New(b.String())
		return nil, err
	}

	return b, nil
}

func (srv *Server) GetCommand(path string) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	res, err := http.Get(srv.Url(path))
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
