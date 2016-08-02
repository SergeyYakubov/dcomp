// Package contains dComp commands that can be executed from command line.
// Every CommandXxxx function that is a member of a cmd struct processes dcomp xxxx command
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

var flHelp bool

var outBuf io.Writer = os.Stdout

func printHelp(f *flag.FlagSet) bool {
	if flHelp {
		f.Usage()
		return true
	} else {
		return false
	}
}

// DoCommand takes command name as a parameter and executes corresponding to this name cmd method
func DoCommand(name string, args []string) error {
	commandName := "Command" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	cmd := new(command)

	methodVal := reflect.ValueOf(cmd).MethodByName(commandName)
	if !methodVal.IsValid() {
		return errors.New("Wrong dcomp command: " + name + "\nType 'dcomp --help'")
	}
	cmd.name = name
	cmd.args = args

	method := methodVal.Interface().(func() error)

	return method()
}

// PrintAllCommands prints all available commands (found wihtin methods of cmd)
func PrintAllCommands() {
	fmt.Fprintln(outBuf, "\nCommands:")
	cmd := new(command)
	CmdType := reflect.TypeOf(cmd)
	for i := 0; i < CmdType.NumMethod(); i++ {
		methodVal := CmdType.Method(i)
		if strings.HasPrefix(methodVal.Name, "Command") {
			method := methodVal.Func.Interface().(func(*command) error)
			cmd.name = strings.ToLower(methodVal.Name)[7:]
			cmd.args = []string{"description"}
			method(cmd)
		}
	}
}
