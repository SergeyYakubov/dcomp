package local

import (
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
	"time"
)

type Resource struct {
	updateStatusCmd func(interface{})
	db              database.Agent
}

type Job struct {
	pid         int
	docker_name string
	status      int
}

func (res *Resource) SubmitJob(job structs.JobInfo) error {
	var localJob Job
	localJob.docker_name = job.Id
	go runScript(job.JobDescription, time.Hour*48)
	res.updateStatusCmd(localJob)
	return nil
}

func (res *Resource) SetUpdateStatusCmd(f func(interface{})) {
	res.updateStatusCmd = f
}

func (res *Resource) SetDb(db database.Agent) {
	res.db = db
}
