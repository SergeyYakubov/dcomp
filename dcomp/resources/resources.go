package resources

import (
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type Resource interface {
	SubmitJob(structs.JobInfo) error
	SetDb(database.Agent)
	GetJob(string) (structs.JobStatus, error)
}
