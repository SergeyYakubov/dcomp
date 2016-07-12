package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func createJobinfoFlags(flagset *flag.FlagSet, flags *structs.JobInfo) {
	flagset.StringVar(&flags.Id, "id", "", "Job Id")
}

func (cmd *command) parseJobinfoFlags(flagset *flag.FlagSet, flags *structs.JobInfo) error {

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	if flags.Id == "" || bson.IsObjectIdHex(flags.Id) {
		return nil
	} else {
		return errors.New("wrong job id format")
	}
}

func (cmd *command) CommandJobinfo() error {

	description := "Show job information"

	if cmd.description(description) {
		return nil
	}

	var flags structs.JobInfo
	flagset := cmd.createFlagset(description, "")
	createJobinfoFlags(flagset, &flags)

	if err := cmd.parseJobinfoFlags(flagset, &flags); err != nil {
		return err
	}

	b, err := Server.GetCommand("jobs" + "/" + flags.Id)

	if err != nil {
		return err
	}

	if b.Len() == 0 {
		fmt.Fprintln(OutBuf, "no jobs found")
		return nil
	}

	decoder := json.NewDecoder(b)
	var jobs []structs.JobInfo
	if err := decoder.Decode(&jobs); err != nil {
		return err
	}

	for i, job := range jobs {
		if flags.Id == "" {
			if i == 0 {
				fmt.Fprintf(OutBuf, "%20s  %20s\n", "Id", "Status")
			}
			job.PrintShort(OutBuf)
		} else {
			job.PrintFull(OutBuf)
		}
	}
	return err
}
