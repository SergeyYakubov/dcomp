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
	{args: []string{"-ncpus", "1", "-script", "aaa", "imagename"}},
	{args: []string{"-ncpus", "10", "-script", "aaa", "imagename"}},
	{args: []string{"-nnodes", "10", "-script", "aaa", "imagename"}},
}

var submitFailingTests = []command{
	{args: []string{"imagename"}},
	{args: []string{}},
	{args: []string{"-script", "aaa"}},
	{args: []string{"-script", "aaa", "-ncpus", "-10", "imagename"}},
	{args: []string{"-script", "aaa", "-ncpus", "0", "-nnodes", "0", "imagename"}},
	{args: []string{"-script", "aaa", "-ncpus", "10", "-nnodes", "10", "imagename"}},
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

var uploadData = []struct {
	localname string
	inipath   string
	destdir   string
	isdir     bool
	result    string
}{
	{"file.txt", ".", "dest", false, "dest/file.txt"},
	{"dir/file.txt", "dir/file.txt", "dest", false, "dest/file.txt"},
	{"dir/file.txt", "dir", "dest", false, "dest/dir/file.txt"},
	{"dir/dir2/file.txt", "dir/dir2", "dest", false, "dest/dir2/file.txt"},
	{"dir", ".", "dest", true, "dest/dir"},
	{"dir/dir2", "dir", "dest", true, "dest/dir/dir2"},
	{"dir/dir2", "dir", "dest", true, "dest/dir/dir2"},
	{"dir/dir2", "dir/dir2", "dest", true, "dest/dir2"},
	{"dir/dir2", "dir/dir2", ".", true, "dir2"},
	{"dir", ".", ".", true, "dir"},
	{".", ".", ".", true, ""},
	{"/dir", "/", ".", true, "dir"},
}

func TestGetUploadName(t *testing.T) {
	for _, test := range uploadData {
		res := getUploadName(test.localname, test.inipath, test.destdir, test.isdir)
		assert.Equal(t, test.result, res)
	}

}
