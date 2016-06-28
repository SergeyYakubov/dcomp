package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Srv struct {
	Host string
	Port int
}

func (srv *Srv) ParseUrl(s string) {
	u, _ := url.Parse(s)
	host, port, _ := net.SplitHostPort(u.Host)
	srv.Host = host
	srv.Port, _ = strconv.Atoi(port)

}

func (srv *Srv) Url(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, "/")
	s = strings.TrimRight(s, "/")
	return fmt.Sprintf("http://%s:%d/%s/", srv.Host, srv.Port, s)
}

func (srv *Srv) PostCommand(path string, data interface{}) (b *bytes.Buffer, err error) {
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
	return b, nil
}
