package resources

import (
	"bytes"
	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

type Resource interface {
	SubmitJob(structs.JobInfo, bool) error
	SetDb(database.Agent)
	GetJobStatus(string) (structs.JobStatus, error)
	DeleteJob(string) error
	PatchJob(string, structs.PatchJob) error
	GetLogs(string, bool) (*bytes.Buffer, error)
}
