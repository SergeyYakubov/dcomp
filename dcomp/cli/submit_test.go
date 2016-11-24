package cli

import (
	"bytes"
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
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

var submitRequests = []struct {
	cmd     command
	answer  string
	code    int
	message string
}{
	{command{args: []string{"-upload", "/etc/passwd:", "-script", "-ncpus", "10", "aaa", "imagename"}}, "578359205e935a20adb39a18", 0, "submit with file"},
	{command{args: []string{"-upload", "/etc/shadow:", "-script", "-ncpus", "10", "aaa", "imagename"}}, "denied", 0, "submit with denied file"},
}

func TestSubmitCommandWithFiles(t *testing.T) {
	t.SkipNow()
	outBuf = new(bytes.Buffer)
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range submitRequests {
		err := test.cmd.CommandSubmit()
		if test.code != 0 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Contains(t, outBuf.(*bytes.Buffer).String(), test.answer, test.message)
		outBuf.(*bytes.Buffer).Reset()
	}
}
