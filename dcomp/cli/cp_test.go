package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cpTests = []struct {
	cmd      command
	answer   cpFlags
	errormsg string
	message  string
}{
	{command{args: []string{"description"}}, cpFlags{}, "", "description"},
	{command{args: []string{}}, cpFlags{}, "wrong job id", "wrong job id"},
	{command{args: []string{"1"}}, cpFlags{}, "wrong job id", "wrong job id"},
	{command{args: []string{"578359205e935a20adb39a18", "."}}, cpFlags{}, "absolute", "relative dest"},
	{command{args: []string{"578359205e935a20adb39a18", "/tmp", "."}},
		cpFlags{Id: "578359205e935a20adb39a18", Source: `/tmp`,
			Dest: "./578359205e935a20adb39a18.tgz", Unpack: false}, "", "copy to current dir"},
	{command{args: []string{"578359205e935a20adb39a18", "/tmp", "/bla"}},
		cpFlags{}, "permission", "copy to nonexisting dir"},
	{command{args: []string{"578359205e935a20adb39a18", "/tmp", "/tmp/file.tgz"}},
		cpFlags{Id: "578359205e935a20adb39a18", Source: `/tmp`,
			Dest: "/tmp/file.tgz", Unpack: false}, "", "copy to /tmp"},
	{command{args: []string{"-u", "578359205e935a20adb39a18", "/tmp", "/tmp/file.tgz"}},
		cpFlags{}, "directory", "copy to /tmp"},
	{command{args: []string{"-u", "578359205e935a20adb39a18", "/tmp", "."}},
		cpFlags{Id: "578359205e935a20adb39a18", Source: `/tmp`,
			Dest: "./", Unpack: true}, "", "copy and unpack to current dir"},
}

func (cmd *command) ParseCpFlags(d string) (cpFlags, error) {
	return cmd.parseCpFlags(d)
}

func TestParseCpFlags(t *testing.T) {

	for _, test := range cpTests {
		flags, err := test.cmd.ParseCpFlags("")
		if err == nil {
			assert.Equal(t, flags, test.answer, test.message)
		} else {
			assert.Contains(t, err.Error(), test.errormsg, test.message)
		}

	}

}
