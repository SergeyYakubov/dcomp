package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var setTests = []struct {
	Test     string
	Err      bool
	Messsage string
}{
	{"aaa:bbb", false, "ok"},
	{"aaa:", true, "empty dest"},
	{"aaa:.", true, "relative dest"},
	{"aaa:./", true, "relative dest 2 "},
	{"aaa:/", true, "root dest"},
}

func TestSetTransferFiles(t *testing.T) {

	for _, test := range setTests {
		var setTests TransferFiles
		err := setTests.Set(test.Test)
		if test.Err {
			assert.NotNil(t, err, test.Messsage)
		} else {
			assert.Nil(t, err, test.Messsage)
		}

	}

}

var outputTests = []struct {
	output  string
	err     bool
	result  JobStatus
	message string
}{
	{"2017-01-19T16:08:22 2017-01-23T21:49:59 COMPLETED", false, JobStatus{StartTime: "2017-01-19 16:08:22 +0000 UTC",
		Status: StatusFinished, EndTime: "2017-01-23 21:49:59 +0000 UTC"}, "finished job"},
	{"2017-01-19T16:08:22 Unknown   RUNNING", false, JobStatus{StartTime: "2017-01-19 16:08:22 +0000 UTC",
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
