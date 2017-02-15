package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var uploadData = []struct {
	localname string
	inipath   string
	destdir   string
	isdir     bool
	result    string
}{
	{"file.txt", ".", "dest", false, "dest"},
	{"dir/file.txt", "dir/file.txt", "dest", false, "dest"},
	{"dir/file.txt", "dir", "dest", false, "dest/file.txt"},
	{"dir/dir2/file.txt", "dir/dir2", "dest", false, "dest/file.txt"},
	{"dir", ".", "dest", true, "dest"},
	{"dir/dir2", "dir", "dest", true, "dest/dir2"},
	{"dir/dir2", "dir", "dest", true, "dest/dir2"},
	{"dir/dir2", "dir/dir2", "dest", true, "dest"},
	{".", ".", "dest", true, "dest"},
	{"/dir", "/", "dest", true, "dest/dir"},
	{"/d/5/aaa/aaa.txt", "/d/5/aaa", "/d/6/ddd", false, "/d/6/ddd/aaa.txt"},
}

func TestGetUploadName(t *testing.T) {
	for _, test := range uploadData {
		res := GetUploadName(test.localname, test.inipath, test.destdir, test.isdir)
		assert.Equal(t, test.result, res)
	}

}
