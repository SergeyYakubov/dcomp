package daemon

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"net/url"

	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
	"time"
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

	db = new(jobdatabase.Mongodb)

	//var dbServer server.Server

	//dbServer.Host = "172.17.0.2"
	//dbServer.Port = 27017

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

type getJobRequest struct {
	path    string
	querry  string
	answer  string
	message string
}

var getJobTests = []getJobRequest{
	{"jobs/578359205e935a20adb39a18", "", "102", "get job info"},
	{"jobs/578359205e935a20adb39a18", "?log=true", "hello", "get log"},
	{"jobs/578359205e935a20adb39a18", "?log=true&compress=true", utils.CompressString("hello"), "get compressed log"},
}

func TestRouteGetJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	setConfiguration()
	var srv server.Server
	ts3 := server.CreateMockServer(&srv)
	defer ts3.Close()
	resources = make(map[string]structs.Resource)
	resources["mock"] = structs.Resource{Server: srv}

	db = new(jobdatabase.Mongodb)

	//	var dbServer server.Server

	//	dbServer.Host = "172.17.0.2"
	//	dbServer.Port = 27017

	db.SetServer(&dbServer)
	db.SetDefaults("daemondbdtest")
	err := db.Connect()
	assert.Nil(t, err)

	defer db.Close()

	s := structs.JobInfo{JobDescription: structs.JobDescription{}, Id: "578359205e935a20adb39a18",
		JobStatus: structs.JobStatus{Status: structs.StatusFinished}, Resource: "mock"}

	db.CreateRecord("578359205e935a20adb39a18", &s)

	for _, test := range getJobTests {

		req, err := http.NewRequest("GET", "http://localhost:8000/"+test.path+"/"+test.querry, nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Contains(t, w.Body.String(), test.answer, test.message)
	}

	db.DeleteRecordByID(getTests[0].job.Id)
}

type FilterTests struct {
	queryString    string
	filteredJobIDs []string
	message        string
}

var inLast30 = utils.TimeToString(time.Now().Add(-time.Hour))

var unfilteredJobs = []structs.JobInfo{
	{JobDescription: structs.JobDescription{},
		JobStatus:  structs.JobStatus{Status: structs.StatusFinished},
		SubmitTime: inLast30,
		Resource:   "mock", Id: "1"},
	{JobDescription: structs.JobDescription{},
		JobStatus:  structs.JobStatus{Status: structs.StatusRunning},
		SubmitTime: "2016-05-01T15:04:05Z",
		Resource:   "mock", Id: "2"},
	{JobDescription: structs.JobDescription{Script: "hello"},
		JobStatus:  structs.JobStatus{Status: structs.StatusRunning},
		SubmitTime: "2017-05-02T15:04:05Z",
		Resource:   "mock", Id: "3"},
}

var filterTests = []FilterTests{
	{``, []string{"1", "2", "3"}, "all jobs"},
	{`finishedOnly=true`, []string{"1"}, "finished jobs"},
	{`notFinishedOnly=true`, []string{"2", "3"}, "not finished jobs"},
	{`keyword=hello`, []string{"3"}, "search by keyword"},
	{`last=30`, []string{"1", "3"}, "last 30 days"},
	{`from=2015-05-01&to=2016-05-02`, []string{"2"}, "search by from/to"},
}

func TestFilterJobs(t *testing.T) {
	for _, test := range filterTests {

		filter, _ := url.ParseQuery(test.queryString)

		jobs := filterJobs(unfilteredJobs, filter)

		assert.Equal(t, len(test.filteredJobIDs), len(jobs), test.message)

		for _, job := range jobs {
			assert.Contains(t, test.filteredJobIDs, job.Id, test.message+" "+job.Id)
		}
	}
}
