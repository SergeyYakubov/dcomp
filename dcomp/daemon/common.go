package daemon

import (
	"net/http"
	"time"

	"bytes"
	"encoding/json"
	"io"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func GetJobFromDatabase(r *http.Request) (structs.JobInfo, error) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	var job structs.JobInfo

	var jobs []structs.JobInfo
	if err := db.GetRecordsByID(jobID, &jobs); err != nil {
		return job, errors.New("cannot retrieve database job info: " + err.Error())
	}

	if len(jobs) != 1 {
		return job, errors.New("cannot fin job: " + jobID)
	}

	return jobs[0], nil

}

func createJWT(job structs.JobInfo, r *http.Request) (token string, err error) {
	srv := resources[job.Resource].DataManager
	val := r.Context().Value("authorizationResponce")

	if val == nil {
		err = errors.New("No authorization context")
		return
	}

	auth, ok := val.(*server.AuthorizationResponce)
	if !ok {
		err = errors.New("Bad authorization context")
		return
	}

	var claim server.JobClaim
	claim.UserName = auth.UserName
	claim.JobInd = job.Id

	var c server.CustomClaims
	c.ExtraClaims = &claim
	c.Duration = time.Hour * 2
	token, err = srv.GetAuth().GenerateToken(&c)
	return
}

func encodeJobFilesTransferInfo(job structs.JobInfo, token string) (b *bytes.Buffer, err error) {
	var s structs.JobFilesTransfer
	s.JobID = job.Id
	s.Srv = resources[job.Resource].DataManager
	s.Token = token

	b = new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&s)
	return
}

func writeJWTToken(w http.ResponseWriter, r *http.Request, job structs.JobInfo) error {
	token, err := createJWT(job, r)
	if err != nil {
		return err
	}

	b, err := encodeJobFilesTransferInfo(job, token)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, b)

	return err

}

func SendJWTToken(w http.ResponseWriter, r *http.Request) {

	job, err := GetJobFromDatabase(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := writeJWTToken(w, r, job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
