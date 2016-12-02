// Package containes structures common between various services. Usually these structures used for HTTP requests
package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"io"
	"path"
	"strings"
)

type jobs interface {
	Check() error
}

type transferFile struct {
	Source string
	Dest   string
}

type TransferFiles []transferFile

func (i *TransferFiles) String() string {
	return fmt.Sprint(*i)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (f *TransferFiles) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		pair := strings.Split(dt, ":")
		if len(pair) != 2 {
			return errors.New("use for <source>:<dest> format for uploading files")
		}
		dest := path.Clean(pair[1])
		if dest == "" || strings.HasPrefix(dest, ".") {
			return errors.New("destination should be absolute")
		}
		if dest == "/" {
			return errors.New("cannot use root destination")
		}
		*f = append(*f, transferFile{Source: path.Clean(pair[0]), Dest: dest})
	}
	return nil
}

// Initial structure filled by user during job submittion
type JobDescription struct {
	ImageName     string
	Script        string
	NCPUs         int
	Local         bool
	FilesToUpload TransferFiles
}

// Job status
const (
	// good codes
	StatusSubmitted         = 101
	StatusRunning           = 102
	StatusFinished          = 103
	StatusCreatingContainer = 104
	StatusStartingContainer = 105
	StatusWaitData          = 106

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

type JobFilesTransfer struct {
	Srv   server.Server
	Token string
	JobID string
}

// Structure with complete job information
type JobInfo struct {
	JobDescription
	JobStatus
	Id       string `bson:"_hex_id"`
	Resource string
}

func (d *JobDescription) NeedData() bool {
	return len(d.FilesToUpload) > 0
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
	StatusSubmitted:         "Submitted",
	StatusRunning:           "Running",
	StatusFinished:          "Finished",
	StatusCreatingContainer: "Creating Docker container",
	StatusStartingContainer: "Starting Docker container",
	StatusSubmissionFailed:  "Submission failed",
	StatusErrorFromResource: "Error from resource",
	StatusWaitData:          "Waiting data",
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
