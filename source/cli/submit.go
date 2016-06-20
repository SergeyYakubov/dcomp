package cli

import (
	"flag"
	"fmt"
	"os"
)

type submitFlags struct {
	ImageName string
	Script    string
	Mode      string
}

func createSubmitFlags(flagset *flag.FlagSet, flags *submitFlags) {
	flagset.StringVar(&flags.Mode, "mode", "", "Docker image")
}

func (cmd *Cmd) parseSubmitFlags(flagset *flag.FlagSet, flags *submitFlags) error {

	flagset.Parse(cmd.args)

	if ShowHelp(flagset) {
		os.Exit(0)
	}

	if flagset.NArg() < 1 {
		return cmd.BadCommandOptions("image name not defined")
	}

	flags.ImageName = flagset.Args()[0]

	fmt.Println(flags.ImageName)

	return nil
}

func (cmd *Cmd) CommandSubmit() error {

	description := "Submit job for distributed computing"

	if cmd.ShowDescription(description) {
		return nil
	}

	var flags submitFlags
	flagset := cmd.Subcmd(description, "IMAGE [COMMAND] [ARG...]")

	createSubmitFlags(flagset, &flags)

	if err := cmd.parseSubmitFlags(flagset, &flags); err != nil {
		return err
	}

	return nil
}
