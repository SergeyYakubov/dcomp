package cli

import (
	"flag"
	"fmt"
	"os"
	"stash.desy.de/scm/dc/common_structs"
)

func createSubmitFlags(flagset *flag.FlagSet, flags *commonStructs.JobDescription) {
	flagset.StringVar(&flags.Script, "script", "", "Job script")
	flagset.IntVar(&flags.NCPUs, "ncpus", 1, "Number of CPUs")
}

func (cmd *Cmd) parseSubmitFlags(flagset *flag.FlagSet, flags *commonStructs.JobDescription) error {

	flagset.Parse(cmd.args)

	if ShowHelp(flagset) {
		os.Exit(0)
	}

	if flagset.NArg() < 1 {
		return cmd.BadCommandOptions("image name not defined")
	}

	flags.ImageName = flagset.Args()[0]

	return flags.Check()
}

func (cmd *Cmd) CommandSubmit() error {

	description := "Submit job for distributed computing"

	if cmd.ShowDescription(description) {
		return nil
	}

	var flags commonStructs.JobDescription
	flagset := cmd.Subcmd(description, "IMAGE [COMMAND] [ARG...]")

	createSubmitFlags(flagset, &flags)

	if err := cmd.parseSubmitFlags(flagset, &flags); err != nil {
		return err
	}

	str, err := Server.PostCommand("jobs", &flags)
	fmt.Fprint(OutBuf, str)
	return err
}
