package docker

import (
	"context"
	"demo-platform/model/db"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"time"
)

type DockerCreateForm struct {
	Name string `form:"name" json:"name" binding:"required"`
	Desc string
	ImageID int64 `form:"ImageID" json:"ImageID" binding:"required"`
	ImageName string `form:"ImageSha1" json:"ImageSha1" binding:"required"`
	ImageSha1 string `form:"ImageSha1" json:"ImageSha1" binding:"required"`
}

func CreateContainer(form *DockerCreateForm, user *db.User) (*string, error) {
	//分配服务器 创建docker
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	_, err = cli.ImagePull(ctx, "docker", types.ImagePullOptions{})

	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	
	name := fmt.Sprintf("%s_%s_%s", user.LoginName, form.Name, now.String())
	
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "centos",
		Cmd:   []string{""},
		Entrypoint: []string{""},
	}, nil, nil, name)

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
	return nil,err
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
	duration := time.Duration(10)
	err = cli.ContainerStop(ctx, containerID, &duration)
	if err != nil || c <= 0 {
		return err
	}
	return nil
}

func StartContainer(containerID string, user *db.User) error {
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
	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{
		CheckpointID:  "",
		CheckpointDir: "",
	})
	if err != nil || c <= 0 {
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