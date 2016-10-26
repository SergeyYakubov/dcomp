package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
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
		Resource: "mock", Id: "678359205e935a20adb39a18"}, structs.StatusRunning, "get single job"},
	{structs.JobInfo{JobDescription: structs.JobDescription{},
		Resource: "aaa", Id: "678359205e935a20adb39a18"}, structs.StatusErrorFromResource, "get single job"},
}

func TestGetJobsFromResources(t *testing.T) {
	setConfiguration()
	var srv server.Server
	ts3 := server.CreateMockServer(&srv)
	defer ts3.Close()
	resources = make(map[string]structs.Resource)
	resources["mock"] = structs.Resource{Server: srv}

	db = new(database.Mongodb)

	var dbServer server.Server

	dbServer.Host = "172.17.0.2"
	dbServer.Port = 27017

	db.SetServer(&dbServer)
	db.SetDefaults("daemondbdtest")
	err := db.Connect()
	assert.Nil(t, err)

	defer db.Close()

	for _, test := range getTests {
		var jobs []structs.JobInfo
		db.CreateRecord(test.job.Id, structs.JobInfo{})
		updateJobsStatusFromResources(&test.job)
		assert.Equal(t, test.answer, test.job.Status, "get status from resources")
		err := db.GetRecordsByID(test.job.Id, &jobs)
		assert.Nil(t, err)
		assert.Equal(t, test.answer, jobs[0].Status, "job not updated")
		db.DeleteRecordByID(test.job.Id)
	}
}
