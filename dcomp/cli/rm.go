package cli

import (
	"errors"
	"flag"
	"os"

	"fmt"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// CommandRm removes job with given id from all places (computation queue, database, etc.)
func (cmd *command) CommandRm() error {

	d := "Delete job files and remove from database"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseRmFlags(d)
	if err != nil {
		return err
	}

	b, status, err := daemon.CommandDelete("jobs" + "/" + flags.Id)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New(b.String())
	}

	fmt.Fprintf(outBuf, "Job removed: %s\n", flags.Id)

	return nil
}

func createRmFlags(flagset *flag.FlagSet, flags *structs.JobInfo) {
	flagset.StringVar(&flags.Id, "id", "", "Job Id")
}

func (cmd *command) parseRmFlags(d string) (structs.JobInfo, error) {

	var flags structs.JobInfo
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS]")

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
