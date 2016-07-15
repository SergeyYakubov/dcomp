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

func createJobinfoFlags(flagset *flag.FlagSet, flags *structs.JobInfo) {
	flagset.StringVar(&flags.Id, "id", "", "Job Id")
}

func (cmd *command) parseJobinfoFlags(d string) (structs.JobInfo, error) {

	var flags structs.JobInfo
	flagset := cmd.createFlagset(d, "[OPTIONS]")
	createJobinfoFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	if flags.Id != "" && !bson.IsObjectIdHex(flags.Id) {
		return flags, errors.New("wrong job id format")
	}

	return flags, nil

}

func (cmd *command) CommandPs() error {

	d := "Show job information"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseJobinfoFlags(d)
	if err != nil {
		return err
	}

	b, err := Server.CommandGet("jobs" + "/" + flags.Id)
	if err != nil {
		return err
	}

	jobs, err := decodeJobs(b)
	if err != nil {
		return err
	}

	printJobs(OutBuf, jobs, flags.Id == "")

	return nil
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
		fmt.Fprintln(OutBuf, "no jobs found")
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
