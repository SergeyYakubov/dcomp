package structs

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Id       string `bson:"_hex_id"`
	Status   int
	Resource string
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

func (d *JobInfo) PrintFull(w io.Writer) {
	fmt.Fprintf(w, "Job %s\n", d.Id)
	fmt.Fprintf(w, "Image name: \t\t %s\n", d.ImageName)
	fmt.Fprintf(w, "Script: \t\t %s\n", d.Script)
	fmt.Fprintf(w, "Number of CPUs: \t %d\n", d.NCPUs)
	fmt.Fprintf(w, "Job status: \t\t %d\n", d.Status)
}

func (d *JobInfo) PrintShort(w io.Writer) {
	fmt.Fprintf(w, "%20s  %20d\n", d.Id, d.Status)
}
