package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"bytes"

	"path/filepath"
	"strings"

	"net/http"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"github.com/sergeyyakubov/dcomp/dcomp/utils"
)

func readJobFilesGetter(bIn *bytes.Buffer) (t structs.JobFilesGetter, err error) {
	decoder := json.NewDecoder(bIn)
	err = decoder.Decode(&t)
	if err != nil || t.Srv.Host == "" || t.Srv.Port == 0 || t.JobID == "" || t.Token == "" {
		err = errors.New("cannot retrieve file upload destination")
	}
	auth := server.NewExternalAuth(t.Token)
	t.Srv.SetAuth(auth)

	return
}

func sendReleaseJobCommand(jobID string) (b *bytes.Buffer, err error) {
	b, _, err = daemon.CommandPost("jobs/"+jobID, nil)
	return
}

func uploadFile(t structs.JobFilesGetter, fileInfo uploadInfo, errchan chan error) {

	f, err := os.Open(fileInfo.Path)
	if err != nil {
		errchan <- err
		return
	}
	defer f.Close()

	un := utils.GetUploadName(fileInfo.Path, fileInfo.Source, fileInfo.Dest, fileInfo.Fi.IsDir())

	_, err = t.Srv.UploadData("jobfile/"+t.JobID+"/", un, f, fileInfo.Fi.Size(), fileInfo.Fi.Mode())
	if err != nil {
		errchan <- err
		return
	}

	errchan <- nil

	return

}

type uploadInfo struct {
	Fi     os.FileInfo
	Path   string
	Source string
	Dest   string
}

func getFilesToUpload(files structs.TransferFiles) (listFiles []uploadInfo, err error) {
	listFiles = make([]uploadInfo, 0)

	for _, pair := range files {
		var scan = func(path string, fi os.FileInfo, err error) (e error) {

			if err != nil {
				return err
			}
			if fi.IsDir() {
				if strings.HasPrefix(fi.Name(), ".") && fi.Name() != "." && fi.Name() != ".." {
					return filepath.SkipDir
				}
			} else {
				if strings.HasPrefix(fi.Name(), ".") {
					return nil
				}

				listFiles = append(listFiles, uploadInfo{fi, path, pair.Source, pair.Dest})

			}
			return nil
		}

		if err = filepath.Walk(pair.Source, scan); err != nil {
			return
		}

	}

	return

}

func uploadFiles(t structs.JobFilesGetter, files structs.TransferFiles) error {

	listFiles, err := getFilesToUpload(files)
	if err != nil {
		return err
	}

	errchan := make(chan error)

	maxParallelRequests := 50

	nrequests := 0
	for _, fileInfo := range listFiles {
		if nrequests < maxParallelRequests {
			go uploadFile(t, fileInfo, errchan)
			nrequests++
		}
		if nrequests == maxParallelRequests {
			err := <-errchan
			nrequests--
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < nrequests; i++ {
		err1 := <-errchan
		if err1 != nil {
			if err == nil {
				err = err1
			} else {
				err = errors.New(err.Error() + err1.Error())
			}
		}
	}
	return err

}

// CommandSubmit sends damon a command to submit a new job. Job id is printed on success, error message otherwise.
func (cmd *command) CommandSubmit() error {

	description := "Submit job for distributed computing"

	if cmd.description(description) {
		return nil
	}

	var flags structs.JobDescription

	flagset := cmd.createDefaultFlagset(description, "[OPTIONS] IMAGE")

	createSubmitFlags(flagset, &flags)

	if err := cmd.parseSubmitFlags(flagset, &flags); err != nil {
		return err
	}

	b, status, err := daemon.CommandPost("jobs", &flags)
	if err != nil {
		return err
	}

	if status != http.StatusCreated && status != http.StatusAccepted {
		return errors.New(b.String())
	}

	if len(flags.FilesToUpload) > 0 {
		t, err := readJobFilesGetter(b)
		if err != nil {
			return err
		}

		err = uploadFiles(t, flags.FilesToUpload)
		if err != nil {
			// file upload failed, delete job from daemon database
			daemon.CommandDelete("jobs" + "/" + t.JobID)
			return err
		}

		b, err = sendReleaseJobCommand(t.JobID)
		if err != nil {
			return err
		}

	}

	decoder := json.NewDecoder(b)
	var t structs.JobInfo
	if err := decoder.Decode(&t); err != nil {
		return err
	}

	fmt.Fprintf(outBuf, "%s\n", t.Id)
	return nil
}

func createSubmitFlags(flagset *flag.FlagSet, flags *structs.JobDescription) {
	flagset.StringVar(&flags.Script, "script", "", "Job script")
	flagset.StringVar(&flags.JobName, "name", "", "Job name")
	flagset.IntVar(&flags.NCPUs, "ncpus", 0, "Number of CPUs")
	flagset.IntVar(&flags.NNodes, "nnodes", 0, "Number of Nodes")
	flagset.StringVar(&flags.Resource, "resource", "", "Force submit to this resource")
	flagset.Var(&flags.FilesToUpload, "upload", "File(s) to upload")
	flagset.Var(&flags.FilesToMount, "mount", "File(s) to mount")

}

func (cmd *command) parseSubmitFlags(flagset *flag.FlagSet, flags *structs.JobDescription) error {

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	if flagset.NArg() < 1 {
		return cmd.errBadOptions("image name not defined")
	}

	flags.ImageName = flagset.Args()[0]
	return flags.Check()
}
