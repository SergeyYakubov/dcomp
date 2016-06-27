package cli

import (
	"../server"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var CommandTests = []struct {
	cmd    Cmd
	answer string
}{
	{Cmd{"submit", []string{"description"}}, "   submit \t\tSubmit job for distributed computing\n"},
	{Cmd{"submit", []string{"-script", "-ncpus", "10", "aaa", "imagename"}}, "OK\n"},
}

var CommandFailingTests = []struct {
	cmd    Cmd
	answer string
}{
	{Cmd{"dummy", []string{"description"}}, "    \t\tSubmit job for distributed computing\n"},
}

func TestCommand(t *testing.T) {
	OutBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&Server, "daemon")
	defer ts.Close()

	for _, test := range CommandFailingTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := Command(test.cmd.name, test.cmd.args)
		assert.NotNil(t, err, "Should be error")

	}
	for _, test := range CommandTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := Command(test.cmd.name, test.cmd.args)
		assert.Nil(t, err, "Should not be error")
		assert.Equal(t, test.answer, OutBuf.(*bytes.Buffer).String(), "")

	}
}

func TestPrintAllCommands(t *testing.T) {
	OutBuf = new(bytes.Buffer)
	PrintAllCommands()
	assert.Contains(t, OutBuf.(*bytes.Buffer).String(), "submit", "all commands mus have submit")
}
