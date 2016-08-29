package local

import (
	"fmt"
	"os"
	"time"

	"io"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type Resource struct {
	db         database.Agent
	werr, wout io.Writer
}

type localJobInfo struct {
	Message string
	Status  int
	Id      string
}

const (
	ContainerCreated  int = 1
	ContainerStarted      = 2
	ContainerFinished     = 3
	ContainerDeleted      = 4
	ContainerError    int = -1
)

func (res *Resource) SubmitJob(job structs.JobInfo) error {
	li := localJobInfo{"", 0, job.Id}
	_, err := res.db.CreateRecord(job.Id, &li)
	if err != nil {
		return err
	}
	go res.runScript(li, job.JobDescription, time.Hour*48)
	return nil
}

func (res *Resource) updateJobInfo(li localJobInfo, status int, message string) {
	li.Status = status
	li.Message = message
	fmt.Fprintln(res.werr, message)
	res.db.PatchRecord(li.Id, li)
}

func createOutputFiles(job structs.JobDescription) (fout, ferr *os.File, err error) {
	foutname := job.WorkDir + `/out.txt`
	fout, err = os.Create(foutname)
	if err != nil {
		return
	}

	ferrname := job.WorkDir + `/err.txt`
	ferr, err = os.Create(ferrname)
	if err != nil {
		fout.Close()
	}

	return
}

func (res *Resource) runScript(li localJobInfo, job structs.JobDescription, d time.Duration) {

	fout, ferr, err := createOutputFiles(job)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res.werr = ferr
	res.wout = fout
	defer fout.Close()
	defer ferr.Close()

	id, err := createContainer(job)
	if err != nil {
		res.updateJobInfo(li, ContainerError, err.Error())
		return
	}

	res.updateJobInfo(li, ContainerCreated, "")

	if err := startContainer(id); err != nil {
		res.updateJobInfo(li, ContainerError, err.Error())
		deleteContainer(id)
		return
	}

	res.updateJobInfo(li, ContainerStarted, "")

	if err := waitFinished(fout, ferr, id, d); err != nil {
		res.updateJobInfo(li, ContainerError, err.Error())
		deleteContainer(id)
		return
	}

	res.updateJobInfo(li, ContainerFinished, "")

	if err := deleteContainer(id); err != nil {
		res.updateJobInfo(li, ContainerError, err.Error())
		return
	}

	res.updateJobInfo(li, ContainerDeleted, "")
}

func (res *Resource) SetDb(db database.Agent) {
	res.db = db
}
