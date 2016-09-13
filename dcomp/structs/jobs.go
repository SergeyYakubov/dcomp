// Package containes structures common between various services. Usually these structures used for HTTP requests
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

// Initial structure filled by user during job submittion
type JobDescription struct {
	ImageName string
	Script    string
	NCPUs     int
	Local     bool
	WorkDir   string
}

// Job status
const (
	// good codes
	StatusSubmitted          = 101
	StatusRunning            = 102
	StatusFinished           = 103
	StatusLoadingDockerImage = 104
	//error codes
	StatusError             = 201
	StatusSubmissionFailed  = 201
	StatusErrorFromResource = 202
)

type JobStatus struct {
	Status    int
	StartTime string
	EndTime   string
	Message   string
}

// Structure with complete job information
type JobInfo struct {
	JobDescription
	JobStatus
	Id       string `bson:"_hex_id"`
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

var jobStatusExplained = map[int]string{
	StatusSubmitted:          "Submitted",
	StatusRunning:            "Running",
	StatusFinished:           "Finished",
	StatusLoadingDockerImage: "Loading Docker image",
	StatusSubmissionFailed:   "Submission failed",
	StatusErrorFromResource:  "Error from resource",
}

func (d *JobInfo) PrintFull(w io.Writer) {
	fmt.Fprintf(w, "%-40s: %s\n", "Job", d.Id)
	fmt.Fprintf(w, "%-40s: %s\n", "Image name", d.ImageName)
	fmt.Fprintf(w, "%-40s: %s\n", "Script", d.Script)
	fmt.Fprintf(w, "%-40s: %d\n", "Number of CPUs", d.NCPUs)
	fmt.Fprintf(w, "%-40s: %s\n", "Allocated resource", d.Resource)
	fmt.Fprintf(w, "%-40s: %s\n", "Status", jobStatusExplained[d.Status])
	if d.Status >= StatusError {
		fmt.Fprintf(w, "%-40s: %s\n", "Message", d.Message)
	}
}

func (d *JobInfo) PrintShort(w io.Writer) {
	fmt.Fprintf(w, "%-40s: %s\n", d.Id, jobStatusExplained[d.Status])
}
