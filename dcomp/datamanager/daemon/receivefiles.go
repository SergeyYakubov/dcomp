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
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"os/user"
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
	if err = binary.Read(r.Body, binary.LittleEndian, &mode); err != nil {
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

	path := filepath.Dir(filename)
	err = utils.MkdirAllWithCh(path, os.ModePerm, uid, gid)

	file, err = os.Create(filename)

	mode, err := getFileMode(r)
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

func routeReceiveJobFile(w http.ResponseWriter, r *http.Request) {

	// get user info, exit if jobID does not match
	auth, err := getUserCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	l, err := io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare reieved file length with info in header
	lexpect, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if l != lexpect {
		http.Error(w, "File length does not match", http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusOK)
}
