package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type testsPS struct {
	job     structs.JobInfo
	answer  int
	message string
}

var getTests = []testsPS{
	{structs.JobInfo{JobDescription: structs.JobDescription{},
		Resource: "mock", Id: "resource"}, structs.StatusFinished, "get single job"},
	{structs.JobInfo{JobDescription: structs.JobDescription{},
		Resource: "aaa", Id: "resource"}, structs.StatusErrorFromResource, "get single job"},
}

func TestGetJobsFromResources(t *testing.T) {
	initialize()
	var srv server.Server
	ts3 := server.CreateMockServer(&srv)
	defer ts3.Close()

	resources["mock"] = structs.Resource{Server: srv}

	for _, test := range getTests {
		updateJobsStatusFromResources(&test.job)
		//		assert.Nil(t, err, test.message)
		assert.Equal(t, test.answer, test.job.Status)
	}
}
