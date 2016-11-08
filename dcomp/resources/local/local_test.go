package local

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"stash.desy.de/scm/dc/main.git/dcomp/database"
	"stash.desy.de/scm/dc/main.git/dcomp/server"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type scriptRequest struct {
	job     structs.JobDescription
	answer  string
	status  int
	message string
}

var runScriptTests = []scriptRequest{
	{structs.JobDescription{ImageName: "centos:7", Script: "echo hi"},
		"hi", structs.StatusFinished, "submit echo script"},
	{structs.JobDescription{ImageName: "centos:7", Script: "hi"},
		"hi", structs.StatusErrorFromResource, "bad command"},
	{structs.JobDescription{ImageName: "max-adm01:0000/nosuchcontainer", Script: "echo hi"},
		"hi", structs.StatusErrorFromResource, "container not exist"},
}

func TestRunScript(t *testing.T) {
	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugintest")
	var res = new(Resource)
	res.SetDb(db)
	db.Connect()
	defer db.Close()

	for _, test := range runScriptTests {
		li := localJobInfo{JobStatus: structs.JobStatus{}, Id: "578359205e935a20adb39a18"}

		res.db.CreateRecord("578359205e935a20adb39a18", &li)
		res.Basedir = "/tmp"
		res.runScript(li, test.job, 1*time.Hour)
		var li_res []localJobInfo
		res.db.GetRecordsByID("578359205e935a20adb39a18", &li_res)
		assert.Equal(t, test.status, li_res[0].Status)
		res.db.DeleteRecordByID("578359205e935a20adb39a18")
	}

}

func TestGetJob(t *testing.T) {
	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugintest")
	var res = new(Resource)
	res.SetDb(db)
	db.Connect()
	defer db.Close()

	id := "578359205e935a20adb39a18"

	li := localJobInfo{JobStatus: structs.JobStatus{Status: structs.StatusRunning}, Id: id}

	res.db.CreateRecord(id, &li)

	status, err := res.GetJob(id)

	assert.Nil(t, err)
	assert.Equal(t, structs.StatusRunning, status.Status)

	res.db.DeleteRecordByID(id)

	status, err = res.GetJob(id)
	assert.NotNil(t, err)

}

func TestDeleteJob(t *testing.T) {
	var dbsrv server.Server
	dbsrv.Host = "172.17.0.2"
	dbsrv.Port = 27017
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugintest")
	var res = new(Resource)
	res.SetDb(db)
	db.Connect()
	defer db.Close()

	id := "578359205e935a20adb39a18"

	job := structs.JobDescription{ImageName: "centos:7", Script: "sleep 100"}

	ji := structs.JobInfo{JobDescription: job, Id: id}

	res.Basedir = "/tmp"
	res.SubmitJob(ji)

	time.Sleep(time.Second * 1)
	err := res.DeleteJob(id)

	assert.Nil(t, err)

	err = res.DeleteJob(id)
	assert.NotNil(t, err)

}
