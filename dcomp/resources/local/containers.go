package local

import (
	"strings"

	"time"

	"io"

	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
	"github.com/dcomp/dcomp/structs"
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

func createContainer(job structs.JobDescription) (string, error) {

	cmd := strings.Fields(job.Script)
	config := container.Config{Image: job.ImageName, AttachStdout: false,
		AttachStderr: false, Cmd: cmd}
	resp, err := cli.ContainerCreate(context.Background(), &config, nil, nil, "")
	if err != nil {
		if client.IsErrImageNotFound(err) {
			options := types.ImageCreateOptions{}
			resp, err := cli.ImageCreate(context.Background(), job.ImageName, options)
			if err != nil {
				return "", err
			}
			defer resp.Close()
			_, err = ioutil.ReadAll(resp)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
		resp, errRetry := cli.ContainerCreate(context.Background(), &config, nil, nil, "")
		if errRetry != nil {
			return "", errRetry
		} else {
			return resp.ID, nil
		}
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

// waitFinished read log files in follow mode, blocking execution until container stops or timeout
func waitFinished(wout io.Writer, id string, d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true,
		Timestamps: false, Details: false}
	reader, err := cli.ContainerLogs(ctx, id, options)

	if err != nil {
		return err
	}

	defer reader.Close()
	_, err = stdcopy.StdCopy(wout, wout, reader)

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
