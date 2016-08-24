package mock

import (
	"errors"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type MockResource struct {
}

func (res *MockResource) SubmitJob(job structs.JobDescription) (interface{}, error) {
	if job.ImageName == "errorsubmit" {
		return nil, errors.New("error submitting job")
	}
	return "12345", nil
}