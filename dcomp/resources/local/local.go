package local

import "stash.desy.de/scm/dc/main.git/dcomp/structs"

type Resource struct {
}

type Job struct {
	pid       int
	docker_id string
}

func (res *Resource) SubmitJob(job structs.JobDescription) (interface{}, error) {
	var localJob Job
	localJob.pid = 123
	localJob.docker_id = "c0e03ec82bc2"
	return localJob, nil
}
