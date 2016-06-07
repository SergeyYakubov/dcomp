package cli

import (
	"fmt"
)

func (cli *Cli) CommandRun(args []string) {
	fmt.Println(args)
	return
}
