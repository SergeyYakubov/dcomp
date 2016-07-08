package structs

import (
	"encoding/json"
	"errors"
	"io"
)

type jobs interface {
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

func (d *JobDescription) Check() error {
	if d.NCPUs <= 0 {
		return errors.New("number of cpus should be > 0")
	}

	if d.ImageName == "" {
		return errors.New("image name should be set")
	}

	if d.Script == "" {
		return errors.New("job script should be set")
	}
	return nil
}

func Decode(r io.Reader, t jobs) bool {

	if r == nil {
		return false
	}

	d := json.NewDecoder(r)

	if d.Decode(t) != nil || t.Check() != nil {
		return false
	}

	return true
}
