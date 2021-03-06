package docker

import (
	"bufio"
	"context"
	"demo-platform/model/db"
	"demo-platform/services/docker"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"testing"
)

func TestDockerAttach(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	c, err := cli.ContainerInspect(ctx, "id")

	if err != nil {
		panic(err)
	}

	if !c.State.Running {
		fmt.Printf("You cannot attach to a stopped container, start it first")
	}

	if c.State.Paused {
		fmt.Printf("You cannot attach to a paused container, unpause it first")
	}

	res, err := cli.ContainerAttach(ctx, "some-nginx", types.ContainerAttachOptions{
		Stream:     false,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		DetachKeys: "",
		Logs:       false,
	})
	defer res.Close()

	if err != nil {
		panic(err)
	}
	_, err = res.Conn.Write([]byte("ls"))
	if err != nil {
		panic(err)
	}
	for ;; {
		data, _, _:= res.Reader.ReadLine()
		fmt.Printf(string(data))
	}
}

func TestDockercli(t *testing.T) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	reader, err := cli.ImagePull(ctx, "centos", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"/bin/sh"},
		Entrypoint: []string{""},
	}, nil, nil, "test_run")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}



	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}

func TestAttach(t *testing.T)  {

}

func TestStdin(t *testing.T) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(input.Text())
	}
}

func TestDockerCreateService(t *testing.T) {
	_ = db.SetupDatabase()

	ID, err := docker.CreateContainer(&docker.CreateForm{
		Name:      "test2",
		Desc:      "",
		ImageID:   1,
		ImageName: "",
		ImageSha1: "",
	}, &db.User{
		Id: 1,
		Name:"asuis",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(ID.ID)
}

func TestDockerStopService(t *testing.T) {
}

func TestDockerStartService(t *testing.T) {
	_ = db.SetupDatabase()

	err := docker.StartContainer("5cd91f9200f5ef64798b5b3b6a94bd5d1de218885821ecffbeb2bbb173a573ce", &db.User{
		Id: 1,
	})

	if err != nil {
		panic(err)
	}
}

func TestPullImage(t *testing.T) {
	_ = db.SetupDatabase()
	_ = docker.PullImage("centos", &db.User{Id:1})
}