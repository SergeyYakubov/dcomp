package cli

import (
	"testing"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/stretchr/testify/assert"
)

var waitTests = []command{
	{args: []string{"description"}},
	{args: []string{"778359205e935a20adb39a18"}},
	{args: []string{"-wait-changes", "778359205e935a20adb39a18"}},
	{args: []string{"-status", "running", "878359205e935a20adb39a18"}},
}

var waitFailingTests = []command{
	{args: []string{}},
	{args: []string{"-status", "bla", "578359205e935a20adb39a20"}},
}

func TestWaitCommand(t *testing.T) {
	ts := server.CreateMockServer(&daemon)
	defer ts.Close()

	for _, test := range waitTests {
		err := test.CommandWait()
		assert.Nil(t, err, "Should not be error")
	}
	for _, test := range waitFailingTests {
		err := test.CommandWait()
		assert.NotNil(t, err, "Should be error")
	}

}
