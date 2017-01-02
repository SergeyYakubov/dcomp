package cli

import (
	"flag"
	"os"

	"io"

	"strings"

	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"gopkg.in/mgo.v2/bson"
	"net/url"
)

type lsFlags struct {
	Recursive bool
	Id        string
	Dir       string
}

func getDataTransferInfo(command string) (t structs.JobFilesTransfer, err error) {

	b, err := daemon.CommandGet(command)
	if err != nil {
		return
	}

	t, err = readJobFilesTransferInfo(b)
	return
}

// CommandLs lists job files
func (cmd *command) CommandLs() error {

	d := "Show job files"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseLsFlags(d)
	if err != nil {
		return err
	}

	cmdstr := "jobfile" + "/" + flags.Id + "/?path=" + url.QueryEscape(flags.Dir) + "&nameonly=true"
	if flags.Recursive {
		cmdstr += "&recursive=true"
	}

	dataTransferInfo, err := getDataTransferInfo(cmdstr)

	b, err := dataTransferInfo.Srv.CommandGet(cmdstr)
	if err != nil {
		return err
	}

	_, err = io.Copy(outBuf, b)
	return err

}

func createLsFlags(flagset *flag.FlagSet, flags *lsFlags) {
	flagset.BoolVar(&flags.Recursive, "R", false, "list subdirectories recursively")
}

func (cmd *command) parseLsFlags(d string) (lsFlags, error) {

	var flags lsFlags
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS] <job ID> [<folder>]")

	createLsFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	flags.Id = flagset.Arg(0)
	flags.Dir = flagset.Arg(1)

	if flags.Id == "" || !bson.IsObjectIdHex(flags.Id) {
		return flags, cmd.errBadOptions("wrong job id format")
	}

	if strings.HasPrefix(flags.Dir, ".") {
		return flags, cmd.errBadOptions("destination should be absolute")
	}

	if flags.Dir == "" {
		flags.Dir = "/"
	}

	return flags, nil

}
