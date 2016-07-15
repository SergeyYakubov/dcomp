package cli

import (
	"errors"
	"flag"
	"os"

	"fmt"
	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

func createRmFlags(flagset *flag.FlagSet, flags *structs.JobInfo) {
	flagset.StringVar(&flags.Id, "id", "", "Job Id")
}

func (cmd *command) parseRmFlags(d string) (structs.JobInfo, error) {

	var flags structs.JobInfo
	flagset := cmd.createFlagset(d, "[OPTIONS]")

	createRmFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	if flags.Id == "" {
		return flags, errors.New("job id not set")

	}
	if !bson.IsObjectIdHex(flags.Id) {
		return flags, errors.New("wrong job id format")
	}

	return flags, nil

}

func (cmd *command) CommandRm() error {

	d := "Cancel job"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseRmFlags(d)
	if err != nil {
		return err
	}

	_, err = Server.CommandDelete("jobs" + "/" + flags.Id)
	if err != nil {
		return err
	}

	fmt.Fprintf(OutBuf, "Job deleted: %s\n", flags.Id)

	return nil
}
