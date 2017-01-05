package cli

import (
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
)

var getJobInfoTests = []struct {
	jobid  string
	answer string
}{
	{"578359205e935a20adb39a18", "578359205e935a20adb39a18"},
	{"578359205e935a20adb39a19", "not found"},
}

func TestGetJobInfo(t *testing.T) {
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range getJobInfoTests {
		job, err := getJobInfo(test.jobid)
		if err == nil {
			assert.Equal(t, job.Id, test.answer, "")
		} else {
			assert.Contains(t, err.Error(), test.answer, "")
		}
	}
}
