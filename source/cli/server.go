package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var Server Srv

type Srv struct {
	host string
	port int
}

func (srv *Srv) Url(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, "/")
	s = strings.TrimRight(s, "/")
	return fmt.Sprintf("http://%s:%d/%s/", srv.host, srv.port, s)
}

func (srv *Srv) PostCommand(path string, data interface{}) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return err
	}
	res, err := http.Post(srv.Url(path), "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(OutBuf, res.Body)
	return nil
}
