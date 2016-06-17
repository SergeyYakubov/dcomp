package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Cli struct {
	//Usage func()
}

var flHelp bool

func ShowDescription(args []string, name, description string) bool {
	if len(args) == 1 && args[0] == "description" {
		fmt.Fprintf(os.Stdout, "   %.20s \t\t%s\n", name, description)
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

func Command(name string, args []string) error {
	commandName := "Command" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	dcompcli := new(Cli)

	methodVal := reflect.ValueOf(dcompcli).MethodByName(commandName)
	if !methodVal.IsValid() {
		return errors.New("Wrong dcomp option: " + name)
	}

	method := methodVal.Interface().(func([]string) error)

	return method(args)
}

func PrintAllCliCommands() {
	dcompcli := new(Cli)
	CliType := reflect.TypeOf(dcompcli)
	for i := 0; i < CliType.NumMethod(); i++ {
		methodVal := CliType.Method(i)
		method := methodVal.Func.Interface().(func(*Cli, []string) error)
		method(dcompcli, []string{"description"})
	}
}

// Subcmd is a subcommand of the main "dcomp" command.
// To see all available subcommands, run "dcomp --help"

func Subcmd(name string, description string) *flag.FlagSet {

	flags := flag.NewFlagSet(name, flag.ExitOnError)
	flags.BoolVar(&flHelp, "help", false, "Print usage")
	flags.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage:\t\ndcomp %s [OPTIONS]", name)
		fmt.Fprintf(os.Stdout, "\n\n%s\n", description)
		flags.PrintDefaults()
	}

	return flags
}
