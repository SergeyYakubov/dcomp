package daemon

import (
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"

	"net/url"
	"path/filepath"

	"encoding/binary"
	"os/user"

	"bytes"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
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

func getUserCredentials(r *http.Request) (resp server.AuthorizationResponce, err error) {
	var jc server.JobClaim

	if err = server.JobClaimFromContext(r, &jc); err != nil {
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	if jobID != jc.JobInd {
		err = errors.New("Not authorized for this job " + jobID)
		return
	}
	return jc.AuthorizationResponce, nil

}

func extractUserInformation(auth server.AuthorizationResponce) (uid, gid int, err error) {
	var u *user.User
	u, err = user.Lookup(auth.UserName)
	if err != nil {
		return
	}

	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		return
	}

	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		return
	}
	return

}

func getFileMode(r *http.Request) (mode os.FileMode, err error) {

	m, err := url.QueryUnescape(r.Header.Get("X-Content-Mode"))
	if err != nil {
		return
	}

	b := bytes.NewBufferString(m)

	if err = binary.Read(b, binary.LittleEndian, &mode); err != nil {
		return
	}
	mode |= utils.GroupRead | utils.GroupWrite | utils.OtherWrite | utils.OtherRead

	if mode&utils.UserExecute != 0 {
		mode |= utils.GroupExecute | utils.OtherExecute
	}

	return

}

func createFile(auth server.AuthorizationResponce, r *http.Request) (file *os.File, err error, errorcode int) {

	filename, err := getFileName(r)
	if err != nil {
		errorcode = http.StatusBadRequest
		return
	}

	uid, gid, err := extractUserInformation(auth)
	if err != nil {
		errorcode = http.StatusBadRequest
		return
	}

	mode, err := getFileMode(r)
	if err != nil {
		errorcode = http.StatusBadRequest
		return
	}

	path := filepath.Dir(filename)
	err = utils.MkdirAllWithCh(path, 0777, uid, gid)
	if err != nil {
		errorcode = http.StatusBadRequest
		return
	}

	file, err = os.Create(filename)
	if err != nil {
		errorcode = http.StatusBadRequest
		return
	}

	err = file.Chmod(mode)
	if err != nil {
		errorcode = http.StatusInternalServerError
		return
	}

	if err = file.Chown(uid, gid); err != nil {
		errorcode = http.StatusInternalServerError
		return
	}

	return
}

func processCopyFile(w http.ResponseWriter, r *http.Request, auth server.AuthorizationResponce, mode string) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	if r.Body == nil {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return

	}

	var fi structs.FileCopyInfo

	if json.NewDecoder(r.Body).Decode(&fi) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	uid, gid, err := extractUserInformation(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	source := settings.Resource.BaseDir + "/" + fi.Source + "/" + fi.SourcePath
	dest := settings.Resource.BaseDir + "/" + jobID + "/" + fi.DestPath

	err = utils.CopyPath(source, dest, mode, uid, gid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func routeReceiveJobFile(w http.ResponseWriter, r *http.Request) {

	// get user info, exit if jobID does not match
	auth, err := getUserCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "mount" {
		processCopyFile(w, r, auth, mode)
		return
	}

	// create file and parent directories, when necessary
	// set file ownership and permissions
	file, err, code := createFile(auth, r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	defer file.Close()

	// copy request content to the file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
