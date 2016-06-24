package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var submitOtherTests = []struct {
	cmd    Cmd
	answer string
}{
	{Cmd{args: []string{"description"}}, "    \t\tSubmit job for distributed computing\n"},
}

var submitTests = []Cmd{
	{args: []string{"-script", "aaa", "imagename"}},
	{args: []string{"-script", "-ncpus", "10", "aaa", "imagename"}},
}

var submitFailingTests = []Cmd{
	{args: []string{"imagename"}},
	{args: []string{}},
	{args: []string{"-script", "aaa"}},
	{args: []string{"-script", "aaa", "-ncpus", "-10", "imagename"}},
}

func TestSubmitCommand(t *testing.T) {
	OutBuf = new(bytes.Buffer)

	ts := PrepareMockServer()
	defer ts.Close()

	for _, test := range submitTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := test.CommandSubmit()
		assert.Nil(t, err, "Should not be error")
		assert.Equal(t, "Job submitted\n", OutBuf.(*bytes.Buffer).String(), "")

	}
	for _, test := range submitFailingTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := test.CommandSubmit()
		assert.NotNil(t, err, "Should be error")

	}
	for _, test := range submitOtherTests {
		OutBuf.(*bytes.Buffer).Reset()
		err := test.cmd.CommandSubmit()
		assert.Nil(t, err, "Should not be error")
		assert.Equal(t, test.answer, OutBuf.(*bytes.Buffer).String(), "")

	}
}
