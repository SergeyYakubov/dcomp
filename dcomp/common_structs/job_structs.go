package commonStructs

import (
	"errors"
)

type JobDescription struct {
	ImageName string
	Script    string
	NCPUs     int
}

func (desc *JobDescription) Check() error {
	if desc.NCPUs <= 0 {
		return errors.New("number of cpus should be > 0")
	}

	if desc.ImageName == "" {
		return errors.New("image name should be set")
	}

	if desc.Script == "" {
		return errors.New("job script should be set")
	}
	return nil
}

type JobInfo struct {
	JobDescription
	Id     string
	Status int
}
