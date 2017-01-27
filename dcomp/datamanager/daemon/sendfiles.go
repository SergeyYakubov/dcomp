package daemon

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"archive/tar"
	"compress/gzip"
	"io"
	"net/url"

	"github.com/gorilla/mux"
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
	if nameonly == "false" {
		recursive = "true"
	}

	basepath := settings.Resource.BaseDir + `/` + jobID + `/`
	path := basepath + relpath

	files, err := getFilesToSend(path, recursive == "true")

	if err != nil {
		errstring := strings.Replace(err.Error(), basepath, "", 1)
		http.Error(w, errstring, http.StatusNotFound)
		return
	}

	if nameonly == "true" {
		w.WriteHeader(http.StatusOK)
		for _, f := range files {
			f = strings.Replace(f, basepath, "", 1)
			fmt.Fprintln(w, f)
		}
		return
	} else {
		err := sendTGZ(w, basepath, files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return
}

func sendTGZ(w http.ResponseWriter, basepath string, files []string) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		if err := addPathToTar(tw, basepath, file); err != nil {
			return err
		}
	}
	return nil
}

func addPathToTar(tw *tar.Writer, basepath, path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	header.Name = strings.Replace(path, basepath, "", 1)

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if !fi.IsDir() {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

	}

	return nil
}
