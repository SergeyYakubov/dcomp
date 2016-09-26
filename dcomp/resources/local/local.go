package local

import (
	"fmt"
	"os"
	"time"

	"io"

	"errors"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type Resource struct {
	db         database.Agent
	werr, wout io.Writer
}

type localJobInfo struct {
	structs.JobStatus
	ContainerId string
	Id          string
}

func (res *Resource) SubmitJob(job structs.JobInfo) error {
	li := localJobInfo{JobStatus: structs.JobStatus{}, Id: job.Id}
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
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		return
	}

	li.ContainerId = id
	res.updateJobInfo(li, structs.StatusLoadingDockerImage, "")

	if err := startContainer(id); err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		deleteContainer(id)
		return
	}

	res.updateJobInfo(li, structs.StatusRunning, "")

	if err := waitFinished(fout, ferr, id, d); err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		deleteContainer(id)
		return
	}

	if err := deleteContainer(id); err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		return
	}

	res.updateJobInfo(li, structs.StatusFinished, "")
}

func (res *Resource) SetDb(db database.Agent) {
	res.db = db
}

func (res *Resource) GetJob(id string) (status structs.JobStatus, err error) {
	var li []localJobInfo
	if err := res.db.GetRecordsByID(id, &li); err != nil {
		return status, err
	}

	if len(li) != 1 {
		return status, errors.New("Database error")
	}
	status = li[0].JobStatus

	return
}
func (res *Resource) DeleteJob(id string) error {

	var li []localJobInfo
	if err := res.db.GetRecordsByID(id, &li); err != nil {
		return err
	}

	if len(li) != 1 {
		return errors.New("Cannot find record in local resource database")
	}

	if err := deleteContainer(li[0].ContainerId); err != nil {
		return err
	}

	if err := res.db.DeleteRecordByID(id); err != nil {
		return err
	}

	return nil
}
