package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"testing"
)

var showJobsTests = []struct {
	cmd    command
	answer string
}{
	{command{args: []string{"description"}}, "information"},
	{command{args: []string{}}, "578359205e935a20adb39a18"},
	{command{args: []string{"-id", "578359205e935a20adb39a18"}}, "578359205e935a20adb39a18"},
	{command{args: []string{"-a"}}, "578359205e935a20adb39a19"},
	{command{args: []string{"-id", "578359205e935a20adb39a19"}}, "not found"},
	{command{args: []string{"-id", "1"}}, "wrong"},
}

func TestCommandPs(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range showJobsTests {
		err := test.cmd.CommandPs()
		if err == nil {
			assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, "")
		} else {
			assert.Contains(t, err.Error(), test.answer, "")
		}

		outBuf.(*bytes.Buffer).Reset()
	}

	daemon.Port = -1
	err := showJobsTests[1].cmd.CommandPs()
	assert.NotNil(t, err, "Should be error")

	ts.Close()
	err = showJobsTests[1].cmd.CommandPs()
	assert.NotNil(t, err, "Should be error")
}
