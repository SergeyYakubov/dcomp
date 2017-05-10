package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var setTests = []struct {
	Test     string
	Ok       bool
	Messsage string
}{
	{"aaa:bbb", true, "ok"},
	{"aaa:", false, "empty dest"},
	{"aaa:.", false, "relative dest"},
	{"aaa:./", false, "relative dest 2 "},
	{"aaa:/", false, "root dest"},
}

func TestSetTransferFiles(t *testing.T) {

	for _, test := range setTests {
		var setTests TransferFiles
		err := setTests.Set(test.Test)
		if !test.Ok {
			assert.NotNil(t, err, test.Messsage)
		} else {
			assert.Nil(t, err, test.Messsage)
		}

	}

}

var setJobTests = []struct {
	Test       string
	Ok         bool
	Source     string
	SourcePath string
	DestPath   string
	Messsage   string
}{
	{"job/aaa:bbb", true, "job", "aaa", "bbb", "ok"},
	{"/job/aaa:bbb", false, "", "", "", "resource missed1 "},
	{"aaa:bbb", false, "", "", "", "resource missed 2"},
	{"job/:sss", false, "", "", "", "source path missed "},
	{"job/t/aaa:d/bbb", true, "job", "t/aaa", "d/bbb", "ok"},
	{"job/aaa:", false, "", "", "", "empty dest"},
	{"job/aaa:.", false, "", "", "", "relative dest"},
	{"job/aaa:./", false, "", "", "", "relative dest 2 "},
	{"job/aaa:/", false, "", "", "", "root dest"},
}

func TestSetJobFiles(t *testing.T) {

	for _, test := range setJobTests {
		var Tests FileCopyInfos
		err := Tests.Set(test.Test)
		if !test.Ok {
			assert.NotNil(t, err, test.Messsage)
		} else {
			assert.Nil(t, err, test.Messsage)
			if err == nil {
				assert.Equal(t, test.DestPath, Tests[0].DestPath, test.Messsage)
				assert.Equal(t, test.SourcePath, Tests[0].SourcePath, test.Messsage)
				assert.Equal(t, test.Source, Tests[0].Source, test.Messsage)
			}
		}

	}

}

var outputTests = []struct {
	output  string
	err     bool
	result  JobStatus
	message string
}{
	{"2017-01-19T16:08:22 2017-01-23T21:49:59 COMPLETED", false, JobStatus{StartTime: "2017-01-19T16:08:22Z",
		Status: StatusFinished, EndTime: "2017-01-23T21:49:59Z"}, "finished job"},
	{"2017-01-19T16:08:22 Unknown   RUNNING", false, JobStatus{StartTime: "2017-01-19T16:08:22Z",
		Status: StatusRunning}, "running job"},
	{"2017-01-19T16:08:22 blabla COMPLETED", true, JobStatus{}, "wrong input (elapsed time)"},
	{"2017-01-23T21:49:59 COMPLETED", true, JobStatus{}, "wrong input (num arguments)"},
}

func TestJobStatus_UpdateFromOutput(t *testing.T) {

	for _, test := range outputTests {
		var jobStatus JobStatus
		err := jobStatus.UpdateFromOutput(test.output)
		if test.err {
			assert.NotNil(t, err, test.message)
		} else {
			assert.Equal(t, test.result, jobStatus, test.message)
			assert.Nil(t, err, test.message)
		}

	}

}
