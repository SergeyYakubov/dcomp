package local

import (
	"fmt"
	"os"
	"time"

	"io"

	"bytes"
	"errors"

	"path/filepath"

	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type Resource struct {
	db      jobdatabase.Agent
	wout    io.Writer
	Basedir string
}

type localJobInfo struct {
	structs.JobStatus
	ContainerId string
	Id          string
}

func (res *Resource) SubmitJob(job structs.JobInfo, checkonly bool) error {
	if checkonly {
		return nil
	}

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

func (res *Resource) logFileName(id string) string {
	return res.Basedir + `/` + id + `/job.log`
}

func (res *Resource) createLogFile(id string, job structs.JobDescription) (flog *os.File, err error) {

	fname := res.logFileName(id)

	path := filepath.Dir(fname)
	if err = os.MkdirAll(path, 0777); err != nil {
		return
	}

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
	basedir := res.Basedir + "/" + li.Id + "/"
	id, err := createContainer(job, basedir)
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

	l, err := res.findJob(li.Id)
	if err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
		return
	}
	// only if container was not already stopped by user
	if l.Status != structs.StatusFinished {
		if err := deleteContainer(id); err != nil {
			res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
			return
		}
	}

	res.updateJobInfo(li, structs.StatusFinished, "")
}

func (res *Resource) SetDb(db jobdatabase.Agent) {
	res.db = db
}

func (res *Resource) GetJobStatus(id string) (status structs.JobStatus, err error) {

	li, err := res.findJob(id)
	if err != nil {
		return status, err
	}

	status = li.JobStatus

	return
}

func (res *Resource) findJob(id string) (li localJobInfo, err error) {
	var listjobs []localJobInfo
	if err := res.db.GetRecordsByID(id, &listjobs); err != nil {
		return li, err
	}

	if len(listjobs) != 1 {
		return li, errors.New("Cannot find record in local resource database")
	}
	return listjobs[0], nil
}

func (res *Resource) PatchJob(id string, patch structs.PatchJob) error {

	li, err := res.findJob(id)
	if err != nil {
		return err
	}

	if li.Status == structs.StatusCreatingContainer || li.Status == structs.StatusStartingContainer {
		return errors.New("Job is being initialized, cannot stop now. Try later.")
	}

	if li.Status == structs.StatusRunning {
		if err := deleteContainer(li.ContainerId); err != nil {
			return err
		}
		res.updateJobInfo(li, structs.StatusFinished, "")
	}

	return nil
}

func (res *Resource) DeleteJob(id string) error {

	li, err := res.findJob(id)
	// return OK if job was not found in database
	if err != nil {
		return nil
	}

	if li.Status == structs.StatusCreatingContainer || li.Status == structs.StatusStartingContainer {
		return errors.New("Job is being initialized, cannot delete now. Try later.")
	}

	if li.Status == structs.StatusRunning {

		return errors.New("job is still running")
	}

	if err := res.db.DeleteRecordByID(id); err != nil {
		return err
	}

	fname := res.logFileName(id)
	path := filepath.Dir(fname)
	return os.RemoveAll(path)
}

func (res *Resource) GetLogs(id string, compressed bool) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	li, err := res.findJob(id)
	if err != nil {
		return nil, err
	}

	fname := res.logFileName(li.Id)

	return utils.ReadFile(fname, compressed)

}
