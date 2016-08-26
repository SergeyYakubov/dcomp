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
