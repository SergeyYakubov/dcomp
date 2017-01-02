package daemon

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"net/url"
)

func getFilesToSend(inipath string, recursive bool) (listFiles []string, err error) {
	listFiles = make([]string, 0)
	var scan = func(path string, fi os.FileInfo, err error) (e error) {

		if err != nil {
			return err
		}
		if fi.IsDir() {
			if strings.HasPrefix(fi.Name(), ".") && fi.Name() != "." && fi.Name() != ".." {
				return filepath.SkipDir
			}
			if path == inipath {
				return nil
			}
			listFiles = append(listFiles, path)
			if !recursive {
				return filepath.SkipDir
			}

		} else {
			if strings.HasPrefix(fi.Name(), ".") {
				return nil
			}
			listFiles = append(listFiles, path)
		}
		return nil
	}

	if err = filepath.Walk(inipath, scan); err != nil {
		return
	}

	return
}

func routeSendJobFile(w http.ResponseWriter, r *http.Request) {

	// get user info, exit if jobID does not match
	_, err := getUserCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	relpath, _ := url.QueryUnescape(r.URL.Query().Get("path"))
	nameonly := r.URL.Query().Get("nameonly")
	recursive := r.URL.Query().Get("recursive")

	path := settings.Resource.BaseDir + `/` + jobID + `/` + relpath

	files, err := getFilesToSend(path, recursive == "true")

	if err != nil {
		errstring := strings.Replace(err.Error(), settings.Resource.BaseDir+`/`+jobID+`/`, "", 1)
		http.Error(w, errstring, http.StatusNotFound)
		return
	}

	if nameonly == "true" {
		w.WriteHeader(http.StatusOK)
		for _, f := range files {
			f = strings.Replace(f, settings.Resource.BaseDir+`/`+jobID+`/`, "", 1)
			fmt.Fprintln(w, f)
		}
		return
	}

	return
}
