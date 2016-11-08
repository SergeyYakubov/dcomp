package resources

import (
	"bytes"
	"github.com/dcomp/dcomp/database"
	"github.com/dcomp/dcomp/structs"
)

type Resource interface {
	SubmitJob(structs.JobInfo) error
	SetDb(database.Agent)
	GetJob(string) (structs.JobStatus, error)
	DeleteJob(string) error
	GetLogs(string, bool) (*bytes.Buffer, error)
}
