package cli

import (
	"errors"
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
)

var flHelp bool

var OutBuf io.Writer = os.Stdout

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

func PrintAllCommands() {
	cmd := new(Cmd)
	CmdType := reflect.TypeOf(cmd)
	for i := 0; i < CmdType.NumMethod(); i++ {
		methodVal := CmdType.Method(i)
		if strings.HasPrefix(methodVal.Name, "Command") {
			method := methodVal.Func.Interface().(func(*Cmd) error)
			cmd.name = strings.ToLower(methodVal.Name)[7:]
			cmd.args = []string{"description"}
			method(cmd)
		}
	}
}
