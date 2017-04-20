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

func GetJobFromDatabase(jobID string) (structs.JobInfo, error) {

	var job structs.JobInfo

	var jobs []structs.JobInfo
	if err := db.GetRecordsByID(jobID, &jobs); err != nil {
		return job, errors.New("cannot retrieve database job info: " + err.Error())
	}

	if len(jobs) != 1 {
		return job, errors.New("cannot find job: " + jobID)
	}

	return jobs[0], nil

}

func GetJobFromRequest(r *http.Request) (structs.JobInfo, error) {

	vars := mux.Vars(r)
	jobID := vars["jobID"]

	return GetJobFromDatabase(jobID)
}

func createJWT(job structs.JobInfo, duration time.Duration) (token string, err error) {
	srv := resources[job.Resource].DataManager

	var claim server.JobClaim
	claim.UserName = job.JobUser
	claim.JobInd = job.Id

	var c server.CustomClaims
	c.ExtraClaims = &claim
	c.Duration = duration
	token, err = srv.GetAuth().GenerateToken(&c)
	return
}

func encodeJobFilesTransferInfo(job structs.JobInfo, token string) (b *bytes.Buffer, err error) {
	var s structs.JobFilesGetter
	s.JobID = job.Id
	s.Srv = resources[job.Resource].DataManager
	s.Token = token

	b = new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&s)
	return
}

func writeJWTToken(w io.Writer, job structs.JobInfo, duration time.Duration) error {
	token, err := createJWT(job, duration)
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

	job, err := GetJobFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := writeJWTToken(w, job, time.Second*30); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return

}

func getUser(r *http.Request) (user string, err error) {
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
	user = auth.UserName
	return
}
