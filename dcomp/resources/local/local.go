package local

import (
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"time"
)

type Resource struct {
	updateStatusCmd func(interface{})
}

type Job struct {
	pid         int
	docker_name string
	status      int
}

func (res *Resource) SubmitJob(job structs.JobInfo) (interface{}, error) {
	var localJob Job
	localJob.docker_name = job.Id
	go runScript(job.JobDescription, time.Hour*48)
	res.updateStatusCmd(localJob)
	return localJob, nil
}

func (res *Resource) SetUpdateStatusCmd(f func(interface{})) {
	res.updateStatusCmd = f
}
