package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"testing"
)

var submitOtherTests = []struct {
	cmd    command
	answer string
}{
	{command{args: []string{"description"}}, "Submit"},
}

var submitTests = []command{
	{args: []string{"-script", "aaa", "imagename"}},
	{args: []string{"-script", "-ncpus", "10", "aaa", "imagename"}},
}

var submitFailingTests = []command{
	{args: []string{"imagename"}},
	{args: []string{}},
	{args: []string{"-script", "aaa"}},
	{args: []string{"-script", "aaa", "-ncpus", "-10", "imagename"}},
}

func TestSubmitCommand(t *testing.T) {
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range submitTests {
		err := test.CommandSubmit()
		assert.Nil(t, err, "Should not be error")
		assert.Equal(t, "578359205e935a20adb39a18\n", outBuf.(*bytes.Buffer).String(), "")
		outBuf.(*bytes.Buffer).Reset()
	}
	for _, test := range submitFailingTests {
		err := test.CommandSubmit()
		assert.NotNil(t, err, "Should be error")
	}
	for _, test := range submitOtherTests {
		err := test.cmd.CommandSubmit()
		assert.Nil(t, err, "Should not be error")
		assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, "")
		outBuf.(*bytes.Buffer).Reset()
	}

	daemon.Port = -1
	err := submitTests[0].CommandSubmit()
	assert.NotNil(t, err, "Should be error")

	ts.Close()
	err = submitTests[0].CommandSubmit()
	assert.NotNil(t, err, "Should be error")
}
