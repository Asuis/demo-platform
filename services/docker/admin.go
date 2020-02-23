package docker

import (
	"context"
	"demo-platform/model/db"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"time"
)

func PullImage(imageName string, user *db.User) error {

	now := time.Now()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	affect, err := db.Engine.Count(db.DockerImage{
		Name:        imageName,
	})
	if err != nil {
		return err
	}
	if affect > 0 {
		return nil
	}
	_, err = db.Engine.Insert(db.DockerImage{
		Sha1:        "",
		UserID:      user.Id,
		Name:        imageName,
		Desc:        "",
		RepoPath:    "",
		Created:     now,
		CreatedUnix: now.Unix(),
		Updated:     now,
		UpdatedUnix: now.Unix(),
	})

	if err != nil {
		return err
	}
	return nil
}

func Images() (*[]types.ImageSummary, error){
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	data, err := cli.ImageList(ctx, types.ImageListOptions{
		All:true,
		Filters: filters.Args{},
	})
	if err != nil {
		return nil, err
	}
	return &data, nil
}