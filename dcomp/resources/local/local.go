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
	db      database.Agent
	wout    io.Writer
	Basedir string
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
	if message != "" {
		fmt.Fprintln(res.wout, message)
	}
	res.db.PatchRecord(li.Id, li)
}

func (res *Resource) createLogFile(id string, job structs.JobDescription) (flog *os.File, err error) {

	fname := res.Basedir + `/` + id + `.log`
	flog, err = os.Create(fname)

	return
}

func (res *Resource) runScript(li localJobInfo, job structs.JobDescription, d time.Duration) {

	fout, err := res.createLogFile(li.Id, job)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res.wout = fout
	defer fout.Close()

	res.updateJobInfo(li, structs.StatusCreatingContainer, "")
	id, err := createContainer(job)
	if err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		return
	}

	li.ContainerId = id
	res.updateJobInfo(li, structs.StatusStartingContainer, "")

	if err := startContainer(id); err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		deleteContainer(id)
		return
	}

	res.updateJobInfo(li, structs.StatusRunning, "")

	if err := waitFinished(fout, id, d); err != nil {
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

	if li[0].Status == structs.StatusRunning {
		if err := deleteContainer(li[0].ContainerId); err != nil {
			return err
		}
	}

	if err := res.db.DeleteRecordByID(id); err != nil {
		return err
	}

	return nil
}
