package daemon

import (
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"

	"net/url"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func getFileName(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]

	disp := r.Header.Get("Content-Disposition")

	_, params, err := mime.ParseMediaType(disp)

	if err != nil {
		return "", err
	}

	shortname, ok := params["filename"]
	if !ok {
		return "", errors.New("cannot extract filename")
	}

	if shortname, err = url.QueryUnescape(shortname); err != nil {
		return "", err
	}

	path := settings.Resource.BaseDir + `/` + jobID + `/`
	filename := path + shortname
	return filename, nil

}

func routeReceiveJobFile(w http.ResponseWriter, r *http.Request) {

	filename, err := getFileName(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := filepath.Dir(filename)

	os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(filename)
	file.Chmod(0666)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	l, err := io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lexpect, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if l != lexpect {
		http.Error(w, "file size does not match", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
