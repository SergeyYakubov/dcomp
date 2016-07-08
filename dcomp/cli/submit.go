package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func createSubmitFlags(flagset *flag.FlagSet, flags *structs.JobDescription) {
	flagset.StringVar(&flags.Script, "script", "", "Job script")
	flagset.IntVar(&flags.NCPUs, "ncpus", 1, "Number of CPUs")
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

func (cmd *command) CommandSubmit() error {

	description := "Submit job for distributed computing"

	if cmd.description(description) {
		return nil
	}

	var flags structs.JobDescription

	flagset := cmd.createFlagset(description, "IMAGE [COMMAND] [ARG...]")

	createSubmitFlags(flagset, &flags)

	if err := cmd.parseSubmitFlags(flagset, &flags); err != nil {
		return err
	}

	b, err := Server.PostCommand("jobs", &flags)

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(b)
	var t structs.JobInfo
	if err := decoder.Decode(&t); err != nil {
		return err
	}

	fmt.Fprintf(OutBuf, "%s\n", t.Id)
	return err
}
