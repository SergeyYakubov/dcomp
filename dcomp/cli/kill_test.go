package cli

import (
	"bytes"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
)

var killTests = []struct {
	cmd    command
	answer string
	msg    string
}{
	{command{args: []string{"description"}}, "Kill", "job description"},
	{command{args: []string{"578359205e935a20adb39a18"}}, "killed", "job killed"},
	{command{args: []string{"578359205e935a20adb39a19"}}, "not found", "job not found"},
	{command{args: []string{"1"}}, "wrong", "wrong job id format"},
}

func TestKillCommand(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range killTests {
		err := test.cmd.CommandKill()
		if err == nil {
			assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, test.msg)
		} else {
			assert.Contains(t, err.Error(), test.answer, test.msg)
		}

		outBuf.(*bytes.Buffer).Reset()
	}

	daemon.Port = -1
	err := killTests[1].cmd.CommandKill()
	assert.NotNil(t, err, "Should be error")

	ts.Close()
	err = killTests[1].cmd.CommandKill()
	assert.NotNil(t, err, "Should be error")
}
