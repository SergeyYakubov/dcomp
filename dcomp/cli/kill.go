package cli

import (
	"errors"
	"os"

	"fmt"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type killFlags struct {
	Id string
}

// CommandKill send kill job command to its resource
func (cmd *command) CommandKill() error {

	d := "Kill job"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseKillFlags(d)
	if err != nil {
		return err
	}

	cmdstr := "jobs" + "/" + flags.Id

	data := structs.PatchJob{Status: structs.StatusFinished}

	b, status, err := daemon.CommandPatch(cmdstr, &data)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusNoContent {
		return errors.New(b.String())
	}
	fmt.Fprintf(outBuf, "Job killed: %s\n", flags.Id)

	return nil

}

func (cmd *command) parseKillFlags(d string) (killFlags, error) {

	var flags killFlags
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS]")

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	flags.Id = flagset.Arg(0)

	if flags.Id == "" {
		return flags, errors.New("job id missed ")
	}

	if flags.Id != "" && !bson.IsObjectIdHex(flags.Id) {
		return flags, errors.New("wrong job id format")
	}

	return flags, nil

}
