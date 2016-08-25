package local

import (
	"strings"

	"time"

	"io"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker_source/pkg/stdcopy"
	"golang.org/x/net/context"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

var cli *client.Client

func init() {
	var err error
	defaultHeaders := map[string]string{"User-Agent": "dComp"}
	cli, err = client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
}

func runScript(job structs.JobDescription, d time.Duration) {

	var wout io.Writer
	var werr io.Writer
	id, err := createContainer(job)
	if err != nil {
		return
	}

	if err := startContainer(id); err != nil {
		return
	}

	if err := bReadLogs(wout, werr, id, d); err != nil {
		return
	}

	deleteContainer(id)

}

func createContainer(job structs.JobDescription) (string, error) {

	cmd := strings.Fields(job.Script)
	config := container.Config{Image: job.ImageName, AttachStdout: false,
		AttachStderr: false, Cmd: cmd}
	resp, err := cli.ContainerCreate(context.Background(), &config, nil, nil, "")
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func deleteContainer(id string) error {
	options := types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}
	return cli.ContainerRemove(context.Background(), id, options)
}

func startContainer(id string) error {
	return cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func waitContainer(id string, d time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return cli.ContainerWait(ctx, id)
}

// bReadLogs read log files in follow mode, blocking execution until container stops or timeout
func bReadLogs(wout io.Writer, werr io.Writer, id string, d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true,
		Timestamps: false, Details: false}
	reader, err := cli.ContainerLogs(ctx, id, options)

	if err != nil {
		return err
	}

	defer reader.Close()
	_, err = stdcopy.StdCopy(wout, werr, reader)

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
