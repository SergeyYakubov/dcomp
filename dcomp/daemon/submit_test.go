package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/jobdatabase"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"github.com/stretchr/testify/assert"
)

var submitRouteTests = []request{
	{structs.JobDescription{},
		"jobs/578359205e935a20adb39a18", "POST", http.StatusCreated, "Create job"},
	{structs.JobDescription{},
		"jobs/578359235e935a21510a2244", "POST", http.StatusAccepted, "Create job"},
}

func TestRouteSubmitReleaseJob(t *testing.T) {
	mux := utils.NewRouter(listRoutes)
	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitRouteTests {

		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		r := resources["local"]
		r.Server = srv
		resources["local"] = r
		resources["mock"] = r

		defer ts2.Close()

		req, err := http.NewRequest(test.cmd, "http://localhost:8000/"+test.path+"/", nil)

		assert.Nil(t, err, "Should not be error")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.answer, w.Code, test.message)
		if w.Code == http.StatusCreated {
			assert.Contains(t, w.Body.String(), "578359205e935a20adb39a18", test.message)
		}

		ts3.Close()
	}
}

type submitRequest struct {
	job     structs.JobDescription
	answer  string
	status  int
	message string
}

var submitTests = []submitRequest{
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: "local"}, "578359205e935a20adb39a18",
		structs.StatusSubmitted, "Create job"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: ""}, "578359205e935a20adb39a18",
		structs.StatusSubmitted, "Create job"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: "local",
		FilesToUpload: structs.TransferFiles{
			{"jhjh", "assd"},
			{"jhjh", "assd"},
		}}, "578359205e935a20adb39a18", structs.StatusWaitData, "Wait for files"},
	{structs.JobDescription{ImageName: "aaa", Script: "bbb", NCPUs: 1, Resource: "local",
		FilesToMount: structs.FileCopyInfos{
			{"jhjh", "assd", "local"},
			{"jhjh", "assd", "local"},
		}}, "578359205e935a20adb39a18", structs.StatusWaitData, "Wait for mount files"},

	{structs.JobDescription{ImageName: "nil", Script: "bbb", NCPUs: 1, Resource: "local"}, "available",
		structs.StatusError, "Create job"},
}

func TestTrySubmitJob(t *testing.T) {

	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range submitTests {
		ts2 := server.CreateMockServer(&estimatorServer)

		var srv server.Server
		ts3 := server.CreateMockServer(&srv)
		resources["local"] = structs.Resource{Server: srv}

		defer ts2.Close()

		if test.job.ImageName == "nil" {
			ts3.Close()
		}

		job, err := trySubmitJob("", test.job)
		if test.job.Script == "nil" || test.job.ImageName == "nil" {
			assert.NotNil(t, err, "Should be error")
			if err != nil {
				assert.Contains(t, err.Error(), test.answer, test.message)
			}
			continue
		}

		assert.Nil(t, err, "Should not be error")

		assert.Equal(t, test.status, job.Status, test.message)
		assert.Contains(t, job.Id, test.answer, test.message)
		ts3.Close()
	}
}

func TestFindResource(t *testing.T) {
	for _, test := range submitTests {

		ts := server.CreateMockServer(&estimatorServer)
		if test.job.Script == "nil" {
			ts.Close()
		}

		prio, err := findResources(test.job)
		if test.job.Script == "nil" {
			assert.NotNil(t, err, "Should be error")
			continue
		}

		assert.Nil(t, err, "Should not be error")
		if test.job.Resource != "" {
			assert.Equal(t, 100, prio[test.job.Resource], test.message)
		} else {
			assert.Equal(t, 10, prio["cloud"], test.message)
		}
		ts.Close()
	}
}

func TestWaitRequests(t *testing.T) {
	errchan := make(chan error)
	n := 3
	f := func(errchach chan error, err error) {
		errchan <- err
	}
	go f(errchan, nil)
	go f(errchan, nil)
	go f(errchan, nil)
	err := waitRequests(n, errchan)
	assert.Nil(t, err, "Should not be error")

	go f(errchan, nil)
	go f(errchan, nil)
	go f(errchan, errors.New("aaa"))
	err = waitRequests(n, errchan)
	assert.NotNil(t, err, "Should  be error")
	assert.Equal(t, err.Error(), "aaa")

}

func TestSubmitSingleCopyDataRequest(t *testing.T) {

	type Tests struct {
		job     structs.JobInfo
		fi      structs.FileCopyInfo
		answer  string
		message string
	}

	var tests = []Tests{
		{structs.JobInfo{JobDescription: structs.JobDescription{},
			Resource: "mock", Id: "578359235e935a21510a2244"}, structs.FileCopyInfo{Source: "", DestPath: "", SourcePath: ""}, "wrong id", "wrong source job id "},
		{structs.JobInfo{JobDescription: structs.JobDescription{},
			Resource: "mock", Id: "578359235e935a21510a2244"}, structs.FileCopyInfo{Source: "578359235e935a21510a2244", DestPath: "", SourcePath: ""}, "", "ok"},
		{structs.JobInfo{JobDescription: structs.JobDescription{},
			Resource: "local", Id: "578359235e935a21510a2244"}, structs.FileCopyInfo{Source: "578359205e935a20adb39a18", DestPath: "", SourcePath: ""}, "another resource", "copy from another resource"},
		//		578359205e935a20adb39a18 returns error from mockserver submit job files
		{structs.JobInfo{JobDescription: structs.JobDescription{},
			Resource: "mock", Id: "578359235e935a21510a2244"}, structs.FileCopyInfo{Source: "578359205e935a20adb39a18", DestPath: "", SourcePath: ""}, "error", "status not created"},
	}

	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range tests {
		ts2 := server.CreateMockServer(&estimatorServer)
		var srv server.Server
		ts3 := server.CreateMockServer(&srv)

		var srvdm server.Server
		ts4 := server.CreateMockServer(&srvdm)
		auth2 := server.NewJWTAuth("aaa")
		srvdm.SetAuth(auth2)
		resources["mock"] = structs.Resource{Server: srv, DataManager: srvdm}

		errchan := make(chan error)
		go submitSingleCopyDataRequest(test.job, test.fi, errchan)
		err := <-errchan
		if test.answer == "" {
			assert.Nil(t, err, test.message)
		} else {
			assert.NotNil(t, err, test.message)
			if err != nil {
				assert.Contains(t, err.Error(), test.answer, test.message)
			}
		}
		ts2.Close()
		ts3.Close()
		ts4.Close()

	}
}

func TestSubmitAfterCopyDataRequest(t *testing.T) {

	type Tests struct {
		job     structs.JobInfo
		answer  string
		message string
	}

	var tests = []Tests{
		{structs.JobInfo{JobDescription: structs.JobDescription{FilesToMount: structs.FileCopyInfos{
			structs.FileCopyInfo{Source: "578359235e935a21510a2244",
				DestPath: "", SourcePath: ""},
			structs.FileCopyInfo{Source: "578359235e935a21510a2244",
				DestPath: "", SourcePath: ""}}},
			Resource: "mock", Id: "578359235e935a21510a2244"},
			"", "ok"},
		{structs.JobInfo{JobDescription: structs.JobDescription{FilesToMount: structs.FileCopyInfos{structs.FileCopyInfo{Source: "578359235e935a21510a2244",
			DestPath: "", SourcePath: ""},
			structs.FileCopyInfo{Source: "578359235e935a21510a224",
				DestPath: "", SourcePath: ""}}},
			Resource: "mock", Id: "578359235e935a21510a2244"},
			"wrong", "wrong id"},
	}

	setConfiguration()
	db = new(jobdatabase.Mockdatabase)
	defer func() { db = nil }()

	for _, test := range tests {
		ts2 := server.CreateMockServer(&estimatorServer)
		var srv server.Server
		ts3 := server.CreateMockServer(&srv)

		var srvdm server.Server
		ts4 := server.CreateMockServer(&srvdm)
		auth2 := server.NewJWTAuth("aaa")
		srvdm.SetAuth(auth2)
		resources["mock"] = structs.Resource{Server: srv, DataManager: srvdm}

		err := submitAfterCopyDataRequest(test.job, false)

		if test.answer == "" {
			assert.Nil(t, err, test.message)
		} else {
			assert.NotNil(t, err, test.message)
			if err != nil {
				assert.Contains(t, err.Error(), test.answer, test.message)
			}
		}
		ts2.Close()
		ts3.Close()
		ts4.Close()

	}
}
