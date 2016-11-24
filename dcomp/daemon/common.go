package daemon

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
