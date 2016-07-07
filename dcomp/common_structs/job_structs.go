package commonStructs

import (
	"encoding/json"
	"errors"
	"io"
)

type jobStruct interface {
	Check() error
}

type JobDescription struct {
	ImageName string
	Script    string
	NCPUs     int
}

type JobInfo struct {
	JobDescription
	Id     string
	Status int
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

func DecodeStruct(r io.Reader, t jobStruct) bool {

	if r == nil {
		return false
	}

	decoder := json.NewDecoder(r)

	if decoder.Decode(&t) != nil || t.Check() != nil {
		return false
	}

	return true
}
