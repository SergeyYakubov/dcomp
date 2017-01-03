package cli

import (
	"flag"
	"os"

	"strings"

	"net/url"

	"github.com/sergeyyakubov/dcomp/dcomp/utils"
	"gopkg.in/mgo.v2/bson"
)

type cpFlags struct {
	Id     string
	Source string
	Dest   string
	Unpack bool
}

// CommandCp downloads job files
func (cmd *command) CommandCp() error {

	d := "Download job files"

	if cmd.description(d) {
		return nil
	}

	flags, err := cmd.parseCpFlags(d)
	if err != nil {
		return err
	}

	cmdstr := "jobfile" + "/" + flags.Id + "/?path=" + url.QueryEscape(flags.Source) + "&nameonly=false"

	dataTransferInfo, err := getDataTransferInfo(cmdstr)

	b, err := dataTransferInfo.Srv.CommandGet(cmdstr)
	if err != nil {
		return err
	}

	if flags.Unpack {
		return utils.WriteUnpackedTGZ(flags.Dest, b)
	} else {
		return utils.WriteFile(flags.Dest, b)
	}
}

func createCpFlags(flagset *flag.FlagSet, flags *cpFlags) {
	flagset.BoolVar(&flags.Unpack, "u", false, "unpack tarball after download")
}

func (cmd *command) parseCpFlags(d string) (cpFlags, error) {

	var flags cpFlags
	flagset := cmd.createDefaultFlagset(d, "[OPTIONS] <job ID> <source> <dest>")

	createCpFlags(flagset, &flags)

	flagset.Parse(cmd.args)

	if printHelp(flagset) {
		os.Exit(0)
	}

	flags.Id = flagset.Arg(0)
	flags.Source = flagset.Arg(1)
	flags.Dest = flagset.Arg(2)

	if flags.Id == "" || !bson.IsObjectIdHex(flags.Id) {
		return flags, cmd.errBadOptions("wrong job id format")
	}

	if strings.HasPrefix(flags.Source, ".") {
		return flags, cmd.errBadOptions("source should be absolute")
	}

	if flags.Source == "" {
		return flags, cmd.errBadOptions("source not set")
	}

	if flags.Dest == "" {
		return flags, cmd.errBadOptions("destination not set")
	}

	stat, err := os.Stat(flags.Dest)
	if err == nil {
		if stat.IsDir() {
			flags.Dest += `/`
			if !flags.Unpack {
				flags.Dest += flags.Id + ".tgz"
			}
		} else {
			return flags, cmd.errBadOptions("file already exists: " + flags.Dest)
		}
	} else {

		if flags.Unpack {
			return flags, cmd.errBadOptions("destination should be a directory when unpacking")
		}

		_, err = os.Create(flags.Dest)
		if err != nil {
			return flags, cmd.errBadOptions(err.Error() + flags.Dest)
		}
		os.Remove(flags.Dest)
	}

	return flags, nil
}
