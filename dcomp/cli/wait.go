package cli

import (
	"errors"
	"flag"
	"os"

	"time"

	"strings"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"gopkg.in/mgo.v2/bson"
)

type waitFlags struct {
	Id           string
	WaitStatus   string
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
	var ini_stat = structs.StatusSubmitted

	anyfinish := false
	if flags.WaitStatus == "" {
		flags.WaitStatus = "COMPLETED"
		anyfinish = true
	}

	flags.WaitStatus = strings.ToUpper(flags.WaitStatus)

	status, message := structs.ExplainStatus(flags.WaitStatus)
	if status == structs.StatusUnknown {
		return errors.New(message + " unknonwn")
	}

	for time.Since(start) < flags.TimeOut {
		jobInfo, err := getJobInfo(flags.Id)
		if err != nil {
			return err
		}

		if flags.StatusChange {
			if jobInfo.Status != ini_stat {
				return nil
			}
		} else {

			if status == jobInfo.Status {
				return nil
			}

			if jobInfo.Status/100 == structs.ErrorCode {
				return errors.New("Exit on error, status: " +
					structs.JobStatusExplained[jobInfo.Status])
			}

			if jobInfo.Status/100 == structs.FinishCode {
				if anyfinish {
					return nil
				}
				return errors.New("Exit on job complete, status: " +
					structs.JobStatusExplained[jobInfo.Status])
			}

		}
		time.Sleep(time.Second)
	}

	return errors.New("Timeout, job status undefined")
}

func createWaitFlags(flagset *flag.FlagSet, flags *waitFlags) {
	flagset.DurationVar(&flags.TimeOut, "timeout", time.Second*60, "Timeout")
	flagset.BoolVar(&flags.StatusChange, "wait-changes", false, "Wait until status changes")
	flagset.StringVar(&flags.WaitStatus, "status", "", "Specify status to wait for")

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
