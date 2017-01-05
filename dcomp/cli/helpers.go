package cli

import (
	"errors"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

// getJobInfo retrieves job info from daemon
func getJobInfo(id string) (structs.JobInfo, error) {

	cmdstr := "jobs" + "/" + id

	b, err := daemon.CommandGet(cmdstr)
	if err != nil {
		return structs.JobInfo{}, err
	}

	// jobs are returned as json string containing []structs.JobInfo
	jobs, err := decodeJobs(b)
	if err != nil {
		return structs.JobInfo{}, err
	}

	if len(jobs) == 0 {
		return structs.JobInfo{}, errors.New("Server returned no jobs")
	}

	return jobs[0], nil

}
