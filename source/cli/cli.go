package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type Cli struct {
	//Usage func()
}

func Command(name string, args ...interface{}) {
	commandName := "Command" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	dcompcli := new(Cli)
	function := reflect.ValueOf(dcompcli).MethodByName(commandName)
	if !function.IsValid() {
		fmt.Println("Wrong dcomp option: " + name)
		return
	}
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	function.Call(in)
}
