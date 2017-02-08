package resources

import (
	"bytes"
	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

type Resource interface {
	SubmitJob(structs.JobInfo, bool) error
	SetDb(jobdatabase.Agent)
	GetJobStatus(string) (structs.JobStatus, error)
	DeleteJob(string) error
	PatchJob(string, structs.PatchJob) error
	GetLogs(string, bool) (*bytes.Buffer, error)
}
