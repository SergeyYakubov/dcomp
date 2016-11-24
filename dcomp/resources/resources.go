package resources

import (
	"bytes"
	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

type Resource interface {
	SubmitJob(structs.JobInfo, bool) error
	SetDb(database.Agent)
	GetJob(string) (structs.JobStatus, error)
	DeleteJob(string) error
	GetLogs(string, bool) (*bytes.Buffer, error)
}
