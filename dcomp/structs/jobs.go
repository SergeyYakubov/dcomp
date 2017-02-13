// Package containes structures common between various services. Usually these structures used for HTTP requests
package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
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

func processDestinationParam(d string) (string, error) {
	res := path.Clean(d)

	if res == "" || strings.HasPrefix(res, ".") {
		return "", errors.New("destination should be absolute")
	}
	if res == "/" {
		return "", errors.New("cannot use root destination")
	}
	return res, nil

}

type copyFile struct {
	SourcePath string
	DestPath   string
	Source     string
}

type CopyFiles []copyFile

func (i *CopyFiles) String() string {
	return fmt.Sprint(*i)
}

func (f *CopyFiles) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		source, dest, err := processFilesParam(dt)
		if err != nil {
			return err
		}

		pair := strings.SplitN(source, "/", 2)
		if len(pair) != 2 {
			return errors.New("Use <source>/<source_path> format to set source")
		}
		source = strings.TrimSpace(pair[0])
		sourcePath := strings.TrimSpace(pair[1])
		if source == "" {
			return errors.New("Empty source")
		}

		*f = append(*f, copyFile{Source: source, SourcePath: sourcePath, DestPath: dest})
	}
	return nil
}

func processFilesParam(d string) (string, string, error) {
	pair := strings.Split(d, ":")
	if len(pair) != 2 {
		return "", "", errors.New("use for <source>:<dest> format")
	}
	source := path.Clean(pair[0])
	dest, err := processDestinationParam(pair[1])
	if err != nil {
		return "", "", err
	}
	return source, dest, nil

}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (f *TransferFiles) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		source, dest, err := processFilesParam(dt)
		if err != nil {
			return err
		}
		*f = append(*f, transferFile{Source: source, Dest: dest})
	}
	return nil
}

// Initial structure filled by user during job submittion
type JobDescription struct {
	ImageName     string
	Script        string
	NCPUs         int
	NNodes        int
	Resource      string
	FilesToUpload TransferFiles
	FilesToMount  CopyFiles
}

type PatchJob struct {
	Status int
}

func (p *PatchJob) Decode(r io.Reader) error {

	if r == nil {
		return errors.New("empty body")
	}

	return json.NewDecoder(r).Decode(p)
}

// Job status
const (
	// finshed codes
	FinishCode      = 1
	StatusCancelled = FinishCode*100 + iota
	StatusFinished
)
const (
	// pending codes
	PendingCode     = 2
	StatusSubmitted = PendingCode*100 + iota
	StatusCreatingContainer
	StatusStartingContainer
	StatusWaitData
	StatusPending
	StatusUserDataCopied
)
const (
	// running codes
	RunningCode   = 3
	StatusRunning = RunningCode*100 + iota
	StatusFinishing
)
const (
	//error codes
	ErrorCode   = 4
	StatusError = ErrorCode*100 + iota
	StatusSubmissionFailed
	StatusErrorFromResource
	StatusFailed
	StatusUnknown
)

type JobStatus struct {
	Status    int
	StartTime string
	EndTime   string
	Message   string
}

func ExplainStatus(statusstr string) (status int, message string) {
	switch statusstr {
	case "COMPLETED":
		status = StatusFinished
	case "CANCELLED":
		status = StatusCancelled
	case "COMPLETING":
		status = StatusFinishing
	case "PENDING":
		status = StatusPending
	case "FAILED":
		status = StatusFailed
	case "TIMEOUT":
		status = StatusFailed
		message = "Job terminated due to timeout"
	case "RUNNING":
		status = StatusRunning
	default:
		status = StatusUnknown
		message = "Status: " + statusstr

	}
	return
}

// UpdateFromOutput updates status by output from an external program
func (s *JobStatus) UpdateFromOutput(status string) error {
	// output has format given by slurm or other programs:
	// elapsed_time status
	// 00:02:36   COMPLETED
	vals := strings.Fields(status)

	if len(vals) != 3 {
		return errors.New("Job not in database " + status)
	}

	timestart := vals[0]
	timeend := vals[1]
	statusstr := vals[2]

	if timestart != "Unknown" {
		ftmstring := "2006-01-02T15:04:05"
		time, err := time.Parse(ftmstring, timestart)
		if err != nil {
			return errors.New("Wrong JobStatus output: " + err.Error())
		}
		s.StartTime = time.String()
	}

	if timeend != "Unknown" {
		ftmstring := "2006-01-02T15:04:05"
		time, err := time.Parse(ftmstring, timeend)
		if err != nil {
			return errors.New("Wrong JobStatus output: " + err.Error())
		}
		s.EndTime = time.String()
	}

	s.Status, s.Message = ExplainStatus(statusstr)

	return nil
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
	JobUser  string
	Id       string `bson:"_hex_id"`
	Resource string
}

func (d *JobDescription) NeedInternalDataCopy() bool {
	return len(d.FilesToMount) > 0
}

func (d *JobDescription) NeedUserDataUpload() bool {
	return len(d.FilesToUpload) > 0
}

func (d *JobDescription) Check() error {
	if d.NCPUs < 0 {
		return errors.New("number of cpus should be > 0")
	}

	if d.NNodes < 0 {
		return errors.New("number of nodes should be > 0")
	}

	if d.NCPUs == 0 && d.NNodes == 0 {
		return errors.New("set number of cpus or number of nodes")
	}

	if d.NCPUs > 0 && d.NNodes > 0 {
		return errors.New("cannot set both number of cpus and number of nodes")
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

var JobStatusExplained = map[int]string{
	StatusSubmitted:         "Submitted",
	StatusRunning:           "Running",
	StatusFinished:          "Finished",
	StatusCancelled:         "Cancelled",
	StatusCreatingContainer: "Creating Docker container",
	StatusStartingContainer: "Starting Docker container",
	StatusSubmissionFailed:  "Submission failed",
	StatusErrorFromResource: "Error from resource",
	StatusWaitData:          "Waiting data",
	StatusPending:           "Pending",
	StatusFailed:            "Failed",
	StatusFinishing:         "Finishing",
}

func (d *JobInfo) PrintFull(w io.Writer) {
	fmt.Fprintf(w, "%-40s: %s\n", "Job", d.Id)
	fmt.Fprintf(w, "%-40s: %s\n", "User", d.JobUser)
	fmt.Fprintf(w, "%-40s: %s\n", "Image name", d.ImageName)
	fmt.Fprintf(w, "%-40s: %s\n", "Script", d.Script)
	fmt.Fprintf(w, "%-40s: %d\n", "Number of CPUs", d.NCPUs)
	fmt.Fprintf(w, "%-40s: %s\n", "Allocated resource", d.Resource)
	fmt.Fprintf(w, "%-40s: %s\n", "Status", JobStatusExplained[d.Status])
	if d.Status >= StatusError {
		fmt.Fprintf(w, "%-40s: %s\n", "Message", d.Message)
	}
}

func (d *JobInfo) PrintShort(w io.Writer) {
	fmt.Fprintf(w, "%-40s: %s\n", d.Id, JobStatusExplained[d.Status])
}
