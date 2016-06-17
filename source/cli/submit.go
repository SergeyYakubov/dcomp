package cli

import (
//	"errors"
//	"fmt"
)

func (cli *Cli) CommandSubmit(args []string) error {

	description := "Submit job for distributed computing"
	name := "submit"

	if ShowDescription(args, name, description) {
		return nil
	}

	flags := Subcmd(name, description)

	var (
		flName = flags.String("image", "", "Docker image")
	)

	flags.Parse(args)

	if ShowHelp(flags) {
		return nil
	}

	if *flName == "" {
		flags.Usage()
		return nil
	}

	return nil
}
