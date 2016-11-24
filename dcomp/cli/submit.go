package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"bytes"

	"github.com/pkg/errors"
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
)

func readJobFilesTransferInfo(bIn *bytes.Buffer) (t structs.JobFilesTransfer, err error) {
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
	return daemon.CommandPost("jobs/"+jobID, nil)
}

func uploadFiles(t structs.JobFilesTransfer, files structs.TransferFiles) error {

	token, _ := t.Srv.GetAuth().GenerateToken(nil)
	fmt.Println("token: ", token)
	return nil
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

	b, err := daemon.CommandPost("jobs", &flags)

	if err != nil {
		return err
	}

	if len(flags.FilesToUpload) > 0 {
		t, err := readJobFilesTransferInfo(b)
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
	flagset.IntVar(&flags.NCPUs, "ncpus", 1, "Number of CPUs")
	flagset.BoolVar(&flags.Local, "local", false, "Submit to local resource")
	flagset.Var(&flags.FilesToUpload, "upload", "File(s) to upload")

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
