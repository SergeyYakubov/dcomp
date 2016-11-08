package local

import (
	"testing"
	"time"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

type request struct {
	job     structs.JobDescription
	answer  string
	message string
}

var submitTests = []request{
	{structs.JobDescription{ImageName: "centos:7", Script: "echo hi"},
		"hi", "submit echo script"},
	{structs.JobDescription{ImageName: "centos:7", Script: "hi"},
		"hi", "bad command"},
	{structs.JobDescription{ImageName: "max-adm01:0000/nosuchcontainer", Script: "echo hi"},
		"hi", "container not exist"},
}

func TestCreateContainer(t *testing.T) {
	for _, test := range submitTests {

		id, err := createContainer(test.job)
		if test.job.ImageName == "max-adm01:0000/nosuchcontainer" {
			assert.NotNil(t, err, test.message)
			continue
		}
		assert.Nil(t, err, "Should not be error")

		deleteContainer(id)
	}
}

func TestStartContainer(t *testing.T) {
	id, err := createContainer(submitTests[0].job)
	assert.Nil(t, err, "Create: should not be error")

	err = startContainer(id)
	assert.Nil(t, err, "Start: not be error")

	err = deleteContainer(id)
	assert.Nil(t, err, "Delete :should not be error")

	id, err = createContainer(submitTests[1].job)
	assert.Nil(t, err, "Create: should not be error")

	err = startContainer(id)
	assert.NotNil(t, err, "Start: should be  be error")

	err = deleteContainer(id)
	assert.Nil(t, err, "Delete :should not be error")

	err = startContainer(id)
	assert.NotNil(t, err, "Second start: Should be error")
}

func TestDeleteContainer(t *testing.T) {
	id, err := createContainer(submitTests[0].job)
	assert.Nil(t, err, "Should not be error")

	err = deleteContainer(id)
	assert.Nil(t, err, "First delete: Should not be error")

	err = deleteContainer(id)
	assert.NotNil(t, err, "Second delete: Should be error")

}

func TestWaitContainer(t *testing.T) {
	job := structs.JobDescription{ImageName: "centos:7", Script: "sleep 10s"}
	id, err := createContainer(job)
	assert.Nil(t, err, "Should not be error")

	err = startContainer(id)
	assert.Nil(t, err, "Start: should not be error")

	t1 := time.Now()
	res, err := waitContainer(id, 10*time.Millisecond)
	assert.NotNil(t, err, "Wait: Should be error")
	assert.Equal(t, -1, res, "Wait: return value should be -1")
	t2 := time.Since(t1)

	if t2.Seconds() > 20*time.Millisecond.Seconds() {
		t.Error("Waited too long")
	}

	err = deleteContainer(id)
	assert.Nil(t, err, "Delete: Should not be error")

	_, err = waitContainer(id, 10*time.Millisecond)
	assert.NotNil(t, err, "Second wait: Should be error")

	job.Script = "sleep 0.1s"
	id, _ = createContainer(job)
	startContainer(id)

	res, err = waitContainer(id, 10*time.Second)
	assert.Nil(t, err, "Wait: Should not be error")
	assert.Equal(t, 0, res, "Wait: return value should be 0")
	err = deleteContainer(id)

}

func TestPrintLogs(t *testing.T) {
	job := structs.JobDescription{ImageName: "centos:7", Script: "echo hi"}
	id, err := createContainer(job)

	err = startContainer(id)

	buf_out := new(bytes.Buffer)
	// echo command, logs are written

	err = waitFinished(buf_out, id, 5*time.Second)
	assert.Equal(t, "hi\n", buf_out.String(), "Ouput should be hi")
	assert.Nil(t, err, "Print logs: should not be error")

	err = deleteContainer(id)

	// long command, exit due to timeout
	job.Script = "sleep 10"
	id, _ = createContainer(job)

	startContainer(id)

	err = waitFinished(buf_out, id, 10*time.Millisecond)
	assert.NotNil(t, err, "Print logs: should be error")

	deleteContainer(id)

}
