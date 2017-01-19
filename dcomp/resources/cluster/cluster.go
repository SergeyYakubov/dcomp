package cluster

import (
	"fmt"

	"io"

	"bytes"
	"errors"

	"io/ioutil"
	"os"

	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"path"

	"github.com/sergeyyakubov/dcomp/dcomp/database"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

type Resource struct {
	db          database.Agent
	wout        io.Writer
	Basedir     string
	TemplateDir string
	Name        string
}

type localJobInfo struct {
	structs.JobStatus
	ClusterJobId string
	Id           string
}

func (res *Resource) executeSubmitCommand(script string) (string, error) {
	f := res.TemplateDir + `/submit.sh`

	cmd := exec.Command(f, script)
	cmd.Dir = path.Dir(script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(err.Error() + " " + string(out))
	}
	words := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(words) != 4 {
		return "", errors.New("Cannot extract slurm job id " + string(out))
	}
	id := words[3]
	return strings.TrimSpace(id), nil
}

func (res *Resource) SubmitJob(job structs.JobInfo, checkonly bool) error {
	if checkonly {
		return nil
	}

	li := localJobInfo{JobStatus: structs.JobStatus{}, Id: job.Id}
	_, err := res.db.CreateRecord(job.Id, &li)
	if err != nil {
		return err
	}

	if err := res.createJobDir(job.Id); err != nil {
		return err
	}

	b, err := res.ProcessSubmitTemplate(job)
	if err != nil {
		return err
	}

	fname := res.jobDir(job.Id) + `/job.sh`
	if err := ioutil.WriteFile(fname, b.Bytes(), 0777); err != nil {
		return err
	}

	li.ClusterJobId, err = res.executeSubmitCommand(fname)

	if err != nil {
		res.updateJobInfo(li, structs.StatusErrorFromResource, err.Error())
	} else {
		res.updateJobInfo(li, structs.StatusSubmitted, "")
	}

	return err
}

func (res *Resource) createMap(job structs.JobInfo) (map[string]string, error) {
	m := make(map[string]string)

	m[`${DCOMP_IMAGE_NAME}`] = job.ImageName
	m[`${DCOMP_SCRIPT}`] = job.Script

	u, err := user.Lookup(job.JobUser)
	if err != nil {
		return m, err
	}

	m[`${DCOMP_UID}`] = u.Uid
	m[`${DCOMP_GID}`] = u.Gid

	if job.NCPUs > 0 {
		m[`${DCOMP_NCPUS}`] = strconv.Itoa(job.NCPUs)
	} else {
		m[`${DCOMP_NCPUS}`] = ""
	}

	if job.NNodes > 0 {
		m[`${DCOMP_NNODES}`] = strconv.Itoa(job.NCPUs)
	} else {
		m[`${DCOMP_NNODES}`] = ""
	}

	m[`${DCOMP_WORKDIR}`] = res.jobDir(job.Id)
	m[`${DCOMP_DOCKER_ARGS}`] = ""

	return m, nil
}

func replaceFromMap(m map[string]string, s string) string {
	for key, val := range m {
		if val == "" {
			if strings.Contains(s, "="+key) {
				return ""
			}
		}
	}
	for key, val := range m {
		s = strings.Replace(s, key, val, -1)
	}

	return s
}

func (res *Resource) ProcessSubmitTemplate(job structs.JobInfo) (b *bytes.Buffer, err error) {

	f := res.TemplateDir + `/batch.sh`
	bt, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}

	lines := strings.Split(string(bt), "\n")

	m, err := res.createMap(job)
	if err != nil {
		return
	}

	b = new(bytes.Buffer)
	for _, line := range lines {
		newline := replaceFromMap(m, line)
		if newline != "" {
			fmt.Fprintln(b, newline)
		}
	}

	return
}

func (res *Resource) DeleteJob(id string) error {

	li, err := res.findJob(id)
	if err != nil {
		return err
	}

	if li.Status == structs.StatusRunning {
	}

	if err := res.db.DeleteRecordByID(id); err != nil {
		return err
	}

	return nil
}

func (res *Resource) GetJobStatus(id string) (status structs.JobStatus, err error) {

	li, err := res.findJob(id)
	if err != nil {
		return status, err
	}

	status = li.JobStatus
	return
}

func (res *Resource) GetLogs(id string, compressed bool) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	li, err := res.findJob(id)
	if err != nil {
		return nil, err
	}

	fname := res.logFileName(li.Id)

	return utils.ReadFile(fname, compressed)
}

func (res *Resource) updateJobInfo(li localJobInfo, status int, message string) {
	li.Status = status
	li.Message = message
	if message != "" {
		//		fmt.Fprintln(res.wout, message)
	}
	res.db.PatchRecord(li.Id, li)
}

func (res *Resource) logFileName(id string) string {
	return res.Basedir + `/` + id + `/job.log`
}

func (res *Resource) jobDir(id string) string {
	return res.Basedir + `/` + id
}

func (res *Resource) createJobDir(id string) error {

	path := res.jobDir(id)
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	return nil
}

func (res *Resource) SetDb(db database.Agent) {
	res.db = db
}

func (res *Resource) findJob(id string) (localJobInfo, error) {
	var listJobs []localJobInfo
	if err := res.db.GetRecordsByID(id, &listJobs); err != nil {
		return localJobInfo{}, err
	}

	if len(listJobs) != 1 {
		return localJobInfo{}, errors.New("Cannot find record in cluster database")
	}
	return listJobs[0], nil
}
