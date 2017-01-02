package cli

import (
	"bytes"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
)

var lsTests = []struct {
	cmd    command
	answer string
}{
	{command{args: []string{"description"}}, "files"},
	{command{args: []string{}}, "wrong job id"},
	{command{args: []string{"1"}}, "wrong job id"},
	{command{args: []string{"578359205e935a20adb39a18", "."}}, "absolute"},
	{command{args: []string{"-R", "578359205e935a20adb39a18", "."}}, "absolute"},
}

func TestCommandLs(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range lsTests {
		err := test.cmd.CommandLs()
		if err == nil {
			assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, "")
		} else {
			assert.Contains(t, err.Error(), test.answer, "")
		}

		outBuf.(*bytes.Buffer).Reset()
	}

	daemon.Port = -1
	err := showJobsTests[1].cmd.CommandLs()
	assert.NotNil(t, err, "Should be error")

	ts.Close()
	err = showJobsTests[1].cmd.CommandLs()
	assert.NotNil(t, err, "Should be error")
}
