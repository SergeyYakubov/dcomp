package cli

import (
	"errors"
	"flag"
	"fmt"
)

type Cmd struct {
	name string
	args []string
}

func (cmd *Cmd) ShowDescription(description string) bool {
	if len(cmd.args) == 1 && cmd.args[0] == "description" {
		fmt.Fprintf(OutBuf, "   %.20s \t\t%s\n", cmd.name, description)
		return true
	}
	return false
}

func (cmd *Cmd) BadCommandOptions(err string) error {
	return errors.New("dcomp " + cmd.name + ": " + err + "\nType 'dcomp " + cmd.name + " --help'")
}

// Subcmd is a subcommand of the main "dcomp" command.
// To see all available subcommands, run "dcomp --help"

func (cmd *Cmd) Subcmd(description, args string) *flag.FlagSet {

	flags := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	flags.BoolVar(&flHelp, "help", false, "Print usage")
	flags.Usage = func() {
		fmt.Fprintf(OutBuf, "Usage:\t\ndcomp %s [OPTIONS] "+args, cmd.name)
		fmt.Fprintf(OutBuf, "\n\n%s\n", description)
		flags.PrintDefaults()
	}

	return flags
}
