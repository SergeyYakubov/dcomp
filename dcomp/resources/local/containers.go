package local

import (
	"strings"

	"time"

	"io"

	"io/ioutil"

	"fmt"

	"os/user"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fsouza/go-dockerclient"
	"github.com/sergeyyakubov/dcomp/dcomp/structs"
	"golang.org/x/net/context"
)

var cli *client.Client

func createTCPClient(host string) *docker.Client {
	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	path := curUser.HomeDir + "/.docker"
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, err := docker.NewTLSClient(host, cert, key, ca)
	if err != nil {
		panic(err)
	}
	return client
}

func InitDockerClient(host string) {
	var err error
	defaultHeaders := map[string]string{"User-Agent": "dComp"}

	if strings.Contains(host, "tcp") {
		cli, err = client.NewClient(host, "v1.24", createTCPClient(host).HTTPClient, defaultHeaders)
	} else {
		cli, err = client.NewClient(host, "v1.24", nil, defaultHeaders)
	}

	if err != nil {
		panic(err)
	}
}

func dockerVolumePair(basedir, dest string) string {
	if !strings.HasPrefix(dest, "/") {
		dest = "/" + dest
	}
	return basedir + dest + ":" + dest

}

func prepareBinds(job structs.JobDescription, basedir string) []string {

	binds := make([]string, 0)
	for _, pair := range job.FilesToUpload {
		binds = append(binds, dockerVolumePair(basedir, pair.Dest))
	}
	for _, pair := range job.FilesToMount {
		binds = append(binds, dockerVolumePair(basedir, pair.DestPath))
	}

	return binds
}

func createContainer(job structs.JobDescription, basedir string) (string, error) {
	var cmd []string
	if job.Script != "" {
		cmd = strings.Fields(job.Script)
	}
	config := container.Config{Image: job.ImageName, AttachStdout: false,
		AttachStderr: false, Cmd: cmd}
	hostConfig := container.HostConfig{Binds: prepareBinds(job, basedir)}
	resp, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, "")
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
