package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"io"

	"bytes"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type psFlags struct {
	Id           string
	ShowFinished bool
}

// CommandPs retrieves jobs info from daemon and prints info
func (cmd *command) CommandPs() error {

	d := "Show job information"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parsePsFlags(d)
	if err != nil {
		return err
	}

	cmdstr := "jobs" + "/"
	if flags.Id == "" && flags.ShowFinished {
		cmdstr += "?finished=true"
	} else {
		cmdstr += flags.Id
	}
	b, err := daemon.CommandGet(cmdstr)
	if err != nil {
		return err
	}

	// jobs are returned as json string containing []structs.JobInfo
	jobs, err := decodeJobs(b)
	if err != nil {
		return err
	}

	printJobs(outBuf, jobs, flags.Id == "")

	return nil
}

func createPsFlags(flagset *flag.FlagSet, flags *psFlags) {
	flagset.StringVar(&flags.Id, "id", "", "Job Id")
	flagset.BoolVar(&flags.ShowFinished, "a", false, "Show all jobs includng finished")
}

func (cmd *command) parsePsFlags(d string) (psFlags, error) {

	var flags psFlags
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS]")

	createPsFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	if flags.Id != "" && !bson.IsObjectIdHex(flags.Id) {
		return flags, errors.New("wrong job id format")
	}

	return flags, nil

}

func decodeJobs(b *bytes.Buffer) ([]structs.JobInfo, error) {
	decoder := json.NewDecoder(b)
	var jobs []structs.JobInfo
	if b.Len() > 0 {
		if err := decoder.Decode(&jobs); err != nil {
			return jobs, err
		}
	}
	return jobs, nil
}

func printJobs(OutBuf io.Writer, jobs []structs.JobInfo, short bool) {
	if len(jobs) == 0 {
		return
	}

	if short {
		fmt.Fprintf(OutBuf, "%20s  %20s\n", "Id", "Status")
	}

	for _, job := range jobs {
		if short {
			job.PrintShort(OutBuf)
		} else {
			job.PrintFull(OutBuf)
		}
	}

}
