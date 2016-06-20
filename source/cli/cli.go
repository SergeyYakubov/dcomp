package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Cmd struct {
	name string
	args []string
}

var flHelp bool

func (cmd *Cmd) ShowDescription(description string) bool {
	if len(cmd.args) == 1 && cmd.args[0] == "description" {
		fmt.Fprintf(os.Stdout, "   %.20s \t\t%s\n", cmd.name, description)
		return true
	}
	return false
}

func ShowHelp(flags *flag.FlagSet) bool {
	if flHelp {
		flags.Usage()
		return true
	} else {
		return false
	}
}

func (cmd *Cmd) BadCommandOptions(err string) error {
	return errors.New("dcomp " + cmd.name + ": " + err + "\nType 'dcomp " + cmd.name + " --help'")
}

func Command(name string, args []string) error {
	commandName := "Command" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	cmd := new(Cmd)

	methodVal := reflect.ValueOf(cmd).MethodByName(commandName)
	if !methodVal.IsValid() {
		return errors.New("Wrong dcomp command: " + name + "\nType 'dcomp --help'")
	}
	cmd.name = name
	cmd.args = args

	method := methodVal.Interface().(func() error)

	return method()
}

func PrintAllCliCommands() {
	cmd := new(Cmd)
	CmdType := reflect.TypeOf(cmd)
	for i := 0; i < CmdType.NumMethod(); i++ {
		methodVal := CmdType.Method(i)
		method := methodVal.Func.Interface().(func(*Cmd) error)
		cmd.name = strings.ToLower(methodVal.Name)[7:]
		cmd.args = []string{"description"}
		method(cmd)
	}
}

// Subcmd is a subcommand of the main "dcomp" command.
// To see all available subcommands, run "dcomp --help"

func (cmd *Cmd) Subcmd(description, args string) *flag.FlagSet {

	flags := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	flags.BoolVar(&flHelp, "help", false, "Print usage")
	flags.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage:\t\ndcomp %s [OPTIONS] "+args, cmd.name)
		fmt.Fprintf(os.Stdout, "\n\n%s\n", description)
		flags.PrintDefaults()
	}

	return flags
}
