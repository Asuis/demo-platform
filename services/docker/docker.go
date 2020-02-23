package docker

import (
	"context"
	"demo-platform/model/db"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"os"
	"path"
	"time"
)

const WorkDir string = "/var/srv/docker/"

type CreateForm struct {
	Name string `form:"Name" json:"Name" binding:"required"`
	Desc string
	ImageID int64 `form:"ImageID" json:"ImageID" binding:"required"`
	ImageName string `form:"ImageName" json:"ImageName" binding:"required"`
	ImageSha1 string `form:"ImageSha1" json:"ImageSha1" binding:"required"`
}

func CreateContainer(form *CreateForm, user *db.User) (*container.ContainerCreateCreatedBody, error) {
	//分配服务器 创建docker
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	
	now := time.Now()
	
	name := fmt.Sprintf("%s-%s-%d", user.Name, form.Name, now.Unix())
	// todo 查找镜像id是否存在
	var image = db.DockerImage{
		ID: form.ImageID,
	}
	has, err := db.Engine.Get(&image)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("image is not fount")
	}

	_, err = cli.ImagePull(ctx, image.Name, types.ImagePullOptions{})

	if err != nil {
		return nil, err
	}
	p := path.Join(WorkDir, user.Name, form.Name)

	err = os.MkdirAll(p, 0766)

	if err != nil {
		return nil, err
	}

	dockerConf := container.Config{
		Hostname:        "",
		Domainname:      "",
		ExposedPorts:    nat.PortSet{
			"80/tcp": {},
		},
		Tty:true,
		Healthcheck:     nil,
		ArgsEscaped:     false,
		Image:           image.Name,
		Volumes:         nil,
		WorkingDir:      p,
		Cmd:   []string{"/bin/bash"},
		Entrypoint: []string{""},
	}
	

	hostConf := container.HostConfig{
		Binds:           nil,
		LogConfig:       container.LogConfig{},
		NetworkMode:     "",
		PortBindings:    nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
		AutoRemove:      false,
		PublishAllPorts: false,
		Resources: container.Resources{
			Memory:               1024*1024*1024,
			NanoCPUs:             1,
			BlkioWeight:          1,
			CPUCount:             1,
		},
	}
	resp, err := cli.ContainerCreate(ctx, &dockerConf, &hostConf, nil, name)

	if err != nil {
		return nil, err
	}
	
	var c = db.DockerContainer{
		ContainerID: resp.ID,
		OwnerID:     user.Id,
		Name:        form.Name,
		Desc:        form.Desc,
		ImageID:     form.ImageID,
		ImageSha1:   form.ImageSha1,
		Status:      db.DockerBuildSuccess,
		Port:        nil,
		ServerID:    0,
		LogPath:     "",
		Created:     now,
		CreatedUnix: now.Unix(),
		Updated:     now,
		UpdatedUnix: now.Unix(),
	}
	
	_, err = db.Engine.Insert(c)
	if err != nil {
		return nil, err
	}
	return &resp,err
}

type UpdateForm struct {
	ContainerID int64
	Resources container.Resources
}
func UpdateContainer(form *UpdateForm, user *db.User) (*container.ContainerUpdateOKBody, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	c := db.DockerContainer{
		ID:          form.ContainerID,
		OwnerID:     user.Id,
	}
	has, err := db.Engine.Get(&c)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("container is not found")
	}
	body, err := cli.ContainerUpdate(ctx, c.ContainerID, container.UpdateConfig{
		Resources: form.Resources,
	})
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func ListContainers(user *db.User, page int, pageSize int, order string) (*[] db.DockerContainer, error) {
	var list [] db.DockerContainer
	err := db.Engine.Where("OwnerID = ?", user.Id).Limit(pageSize, page).OrderBy(order).Find(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

/**
todo
*/
func StopContainer(containerID string, user *db.User) error {
	c := db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	}
	has, err := db.Engine.Get(&c)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("container not exist")
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	duration := time.Duration(10)
	err = cli.ContainerStop(ctx, containerID, &duration)
	if err != nil {
		return err
	}

	c.Status = db.DockerStoped

	_, err = db.Engine.Update(c)
	if err != nil {
		return err
	}

	return nil
}

func StartContainer(containerID string, user *db.User) error {
	c := db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	}
	has, err := db.Engine.Get(&c)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("container not exist")
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{
	})
	if err != nil {
		return err
	}
	c.Status = db.DockerRunning
	_, err = db.Engine.Update(c)
	if err != nil {
		return err
	}
	return nil
}

func RmContainer(containerID string, user *db.User) error {
	c, err := db.Engine.Count(&db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	})
	if err != nil {
		return err
	}
	if c<=0 {
		return fmt.Errorf("container not exist")
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         false,
	})
	if err != nil {
		return err
	}
	
	_, err = db.Engine.Delete(&db.DockerContainer{
		ContainerID:containerID,
		OwnerID: user.Id,
	})

	if err != nil{
		return err
	}
	
	return nil
}

func StatusContainer(containerID string, user *db.User) (*types.ContainerStats, error) {
	c, err := db.Engine.Count(&db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	})

	if err != nil {
		return nil, err
	}
	if c<=0 {
		return nil, fmt.Errorf("container not exist")
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	resp, err := cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func InfoContainer(containerID string, user *db.User) (*db.DockerContainer, error) {
	var c = db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	}
	has, err := db.Engine.Get(&c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("container not exist")
	}
	return &c, nil
}

func RestartContainer(containerID string, user *db.User) error {
	c, err := db.Engine.Count(&db.DockerContainer{
		ContainerID: containerID,
		OwnerID:user.Id,
	})
	if err != nil {
		return err
	}
	if c<=0 {
		return fmt.Errorf("container not exist")
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	//todo
	dur := time.Duration(60)

	err = cli.ContainerRestart(ctx, containerID, &dur)

	if err != nil {
		return err
	}

	return nil
}

func Containers() (*[]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	data, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Latest:  false,
	})
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func ContainerInfo(ID string) (*types.ContainerJSON, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	data, err := cli.ContainerInspect(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func ContainerBuild(repoUrl string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	cli.ImageBuild(ctx, nil, types.ImageBuildOptions{
		Memory:         1024*1024*1024,
		Dockerfile:     "",
	})
}