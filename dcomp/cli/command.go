package cli

import (
	"errors"
	"flag"
	"fmt"
)

type command struct {
	name string
	args []string
}

func (cmd *command) description(d string) bool {
	if len(cmd.args) == 1 && cmd.args[0] == "description" {
		fmt.Fprintf(OutBuf, "   %-10s %s\n", cmd.name, d)
		return true
	}
	return false
}

func (cmd *command) errBadOptions(err string) error {
	return errors.New("dcomp " + cmd.name + ": " + err + "\nType 'dcomp " + cmd.name + " --help'")
}

// Subcmd is a subcommand of the main "dcomp" command.
// To see all available subcommands, run "dcomp --help"

func (cmd *command) createFlagset(description, args string) *flag.FlagSet {

	flags := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	flags.BoolVar(&flHelp, "help", false, "Print usage")
	flags.Usage = func() {
		fmt.Fprintf(OutBuf, "Usage:\t\ndcomp %s "+args, cmd.name)
		fmt.Fprintf(OutBuf, "\n\n%s\n", description)
		flags.PrintDefaults()
	}

	return flags
}
