package cluster

import (
	"testing"

	"os/user"
	"strconv"

	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

type scriptRequest struct {
	job     structs.JobDescription
	answer  string
	status  int
	message string
}

type config struct {
	BaseDir     string
	TemplateDir string
}

var runScriptTests = []scriptRequest{
	{structs.JobDescription{ImageName: "centos:7", Script: "echo hi", NCPUs: 1},
		"hi", structs.StatusFinished, "submit echo script"},
	{structs.JobDescription{ImageName: "centos:7", Script: "echo hi"},
		"hi", structs.StatusFinished, "submit echo script"},
}

func TestPrepareScript(t *testing.T) {
	db := createdb()
	var res = new(Resource)
	res.SetDb(db)
	db.Connect()
	defer db.Close()

	setConfiguration(res)

	id := "578359205e935a20adb39a18"

	for _, test := range runScriptTests {
		li := localJobInfo{JobStatus: structs.JobStatus{}, Id: id}

		res.db.CreateRecord(id, &li)

		var li_res []localJobInfo
		u, _ := user.Current()
		ji := structs.JobInfo{JobDescription: test.job, Id: id, JobUser: u.Username}
		b, err := res.ProcessSubmitTemplate(ji)
		assert.Nil(t, err, "")
		assert.Contains(t, b.String(), `--workdir=/dcompdata/578359205e935a20adb39a18`, test.message)
		if test.job.NCPUs > 0 {
			assert.Contains(t, b.String(), strconv.Itoa(test.job.NCPUs), test.message)
		} else {
			assert.NotContains(t, b.String(), "ntasks", test.message)
		}
		assert.Contains(t, b.String(), test.job.ImageName, test.message)
		assert.Contains(t, b.String(), test.job.Script, test.message)
		assert.Contains(t, b.String(), u.Gid, test.message)
		assert.Contains(t, b.String(), u.Uid, test.message)

		res.db.GetRecordsByID(id, &li_res)
		li_res[0].Status = structs.StatusFinished
		assert.Equal(t, test.status, li_res[0].Status)
		res.db.DeleteRecordByID(id)
	}

}

func setConfiguration(res *Resource) {

	var c config

	utils.ReadYaml(`/etc/dcomp/plugins/slurm/slurm.yaml`, &c)

	res.TemplateDir = c.TemplateDir
	res.Basedir = c.BaseDir
}

func createdb() *database.Mongodb {
	var dbsrv server.Server
	dbsrv.Host = "localhost"
	dbsrv.Port = 27017
	db := new(database.Mongodb)
	db.SetServer(&dbsrv)
	db.SetDefaults("localplugintest")
	return db

}

func TestSubmitJob_Checkonly(t *testing.T) {

	var res = new(Resource)

	id := "578359205e935a20adb39a18"

	job := structs.JobDescription{ImageName: "centos:7", Script: "sleep 100"}

	ji := structs.JobInfo{JobDescription: job, Id: id}
	err := res.SubmitJob(ji, true)
	assert.Nil(t, err)

}
