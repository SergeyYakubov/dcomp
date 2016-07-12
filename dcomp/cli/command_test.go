package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"testing"
)

var CommandTests = []struct {
	cmd    command
	answer string
}{
	{command{"submit", []string{"description"}}, "   submit \t\tSubmit job for distributed computing\n"},
	{command{"submit", []string{"-script", "-ncpus", "10", "aaa", "imagename"}}, "1\n"},
}

var CommandFailingTests = []struct {
	cmd    command
	answer string
}{
	{command{"dummy", []string{"description"}}, "    \t\tSubmit job for distributed computing\n"},
}

func TestCommand(t *testing.T) {
	OutBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&Server)
	defer ts.Close()

	for _, test := range CommandFailingTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := DoCommand(test.cmd.name, test.cmd.args)
		assert.NotNil(t, err, "Should be error")

	}
	for _, test := range CommandTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := DoCommand(test.cmd.name, test.cmd.args)
		assert.Nil(t, err, "Should not be error")
		assert.Equal(t, test.answer, OutBuf.(*bytes.Buffer).String(), "")

	}
}

func TestPrintAllCommands(t *testing.T) {
	OutBuf = new(bytes.Buffer)
	PrintAllCommands()
	assert.Contains(t, OutBuf.(*bytes.Buffer).String(), "submit", "all commands mus have submit")
}
