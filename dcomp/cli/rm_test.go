package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"testing"
)

var rmTests = []struct {
	cmd    command
	answer string
	msg    string
}{
	{command{args: []string{"description"}}, "Cancel", "job description"},
	{command{args: []string{"-id", "578359205e935a20adb39a18"}}, "", "job removed"},
	{command{args: []string{"-id", "578359205e935a20adb39a19"}}, "not found", "job not found"},
	{command{args: []string{"-id", "1"}}, "wrong", "wrong path"},
}

func TestRmCommand(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range rmTests {
		err := test.cmd.CommandRm()
		if err == nil {
			assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, test.msg)
		} else {
			assert.Contains(t, err.Error(), test.answer, test.msg)
		}

		outBuf.(*bytes.Buffer).Reset()
	}

	daemon.Port = -1
	err := rmTests[1].cmd.CommandRm()
	assert.NotNil(t, err, "Should be error")

	ts.Close()
	err = rmTests[1].cmd.CommandRm()
	assert.NotNil(t, err, "Should be error")
}
