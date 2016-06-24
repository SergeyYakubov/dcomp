package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var CommandTests = []struct {
	cmd    Cmd
	answer string
}{
	{Cmd{"submit", []string{"description"}}, "   submit \t\tSubmit job for distributed computing\n"},
	{Cmd{"submit", []string{"-script", "-ncpus", "10", "aaa", "imagename"}}, "Job submitted\n"},
}

var CommandFailingTests = []struct {
	cmd    Cmd
	answer string
}{
	{Cmd{"dummy", []string{"description"}}, "    \t\tSubmit job for distributed computing\n"},
}

func TestCommand(t *testing.T) {
	OutBuf = new(bytes.Buffer)

	ts := PrepareMockServer()
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
