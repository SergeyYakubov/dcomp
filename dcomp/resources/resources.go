package resources

import "stash.desy.de/scm/dc/main.git/dcomp/structs"

type Resource interface {
	SubmitJob(structs.JobDescription) (interface{}, error)
}
