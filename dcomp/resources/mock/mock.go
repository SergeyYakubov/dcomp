package mock

import (
	"errors"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type MockResource struct {
}

func (res *MockResource) SubmitJob(job structs.JobInfo) error {
	if job.ImageName == "errorsubmit" {
		return errors.New("error submitting job")
	}
	return nil
}

func (res *MockResource) SetDb(database.Agent) {

}

func (res *MockResource) GetJob(id string) (status structs.JobStatus, err error) {
	if id == "578359205e935a20adb39a18" {
		status.Status = structs.StatusSubmitted
		return
	}
	err = errors.New("Job not found")
	return
}
