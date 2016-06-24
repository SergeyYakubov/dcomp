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
		return errors.New("Number of cpus should be > 0")
	}

	if desc.ImageName == "" {
		return errors.New("Image name should be set")
	}

	if desc.Script == "" {
		return errors.New("Job script should be set")
	}
	return nil
}
