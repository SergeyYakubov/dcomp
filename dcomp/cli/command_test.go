package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/dcomp/dcomp/server"
	"testing"
)

var CommandTests = []struct {
	cmd    command
	answer string
}{
	{command{"submit", []string{"description"}}, "Submit"},
	{command{"submit", []string{"-script", "-ncpus", "10", "aaa", "imagename"}}, "578359205e935a20adb39a18\n"},
}

var CommandFailingTests = []struct {
	cmd    command
	answer string
}{
	{command{"dummy", []string{"description"}}, "Submit"},
}

func TestCommand(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range CommandFailingTests {
		outBuf.(*bytes.Buffer).Reset()
		err := DoCommand(test.cmd.name, test.cmd.args)
		assert.NotNil(t, err, "Should be error")

	}
	for _, test := range CommandTests {
		outBuf.(*bytes.Buffer).Reset()
		err := DoCommand(test.cmd.name, test.cmd.args)
		assert.Nil(t, err, "Should not be error")
		assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, "")

	}
}

func TestPrintAllCommands(t *testing.T) {
	outBuf = new(bytes.Buffer)
	PrintAllCommands()
	assert.Contains(t, outBuf.(*bytes.Buffer).String(), "submit", "all commands mus have submit")
}
