package mock

import (
	"bytes"
	"compress/gzip"
	"errors"

	"github.com/dcomp/dcomp/database"
	"github.com/dcomp/dcomp/structs"
)

type MockResource struct {
}

func (res *MockResource) SubmitJob(job structs.JobInfo) error {
	if job.ImageName == "errorsubmit" {
		return errors.New("error submitting job")
	}
	return nil
}

func (res *MockResource) SetDb(database.Agent) {

}

func (res *MockResource) GetJob(id string) (status structs.JobStatus, err error) {
	if id == "578359205e935a20adb39a18" {
		status.Status = structs.StatusFinished
		return
	}
	err = errors.New("Job not found")
	return
}

func (res *MockResource) DeleteJob(id string) error {
	if id == "578359205e935a20adb39a18" {
		return nil
	}
	return errors.New("Job not found")
}

func (res *MockResource) GetLogs(id string, compressed bool) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)
	if compressed {
		gz := gzip.NewWriter(b)
		defer gz.Close()
		gz.Write([]byte("hello"))
	} else {
		b.WriteString("hello")
	}
	return b, nil
}
