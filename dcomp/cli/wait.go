package cli

import (
	"errors"
	"flag"
	"os"

	"time"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"gopkg.in/mgo.v2/bson"
)

type waitFlags struct {
	Id           string
	StatusChange bool
	TimeOut      time.Duration
}

// CommandPs retrieves jobs info from daemon and prints info
func (cmd *command) CommandWait() error {

	d := "Wait untill job finishes"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseWaitFlags(d)
	if err != nil {
		return err
	}

	start := time.Now()
	first := true
	var ini_stat int
	for time.Since(start) < flags.TimeOut {
		jobInfo, err := getJobInfo(flags.Id)
		if flags.StatusChange && first {
			first = false
			ini_stat = jobInfo.Status
		}

		if err != nil {
			return err
		}

		if flags.StatusChange {
			if ini_stat != jobInfo.Status {
				return nil
			}
		} else {
			if jobInfo.Status == structs.StatusFinished || jobInfo.Status/100 == structs.ErrorCode {
				return nil
			}
		}
		time.Sleep(time.Second)

	}

	return errors.New("Timeout, job status undefined")
}

func createWaitFlags(flagset *flag.FlagSet, flags *waitFlags) {
	flagset.DurationVar(&flags.TimeOut, "timeout", time.Second*60, "Timeout")
	flagset.BoolVar(&flags.StatusChange, "wait-changes", false, "Wait until status changes")

}

func (cmd *command) parseWaitFlags(d string) (waitFlags, error) {

	var flags waitFlags
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS] <Job ID>")

	createWaitFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	flags.Id = flagset.Arg(0)

	if flags.Id == "" {
		return flags, errors.New("job id missed ")
	}

	if !bson.IsObjectIdHex(flags.Id) {
		return flags, errors.New("wrong job id format: " + flags.Id)
	}
	return flags, nil
}
