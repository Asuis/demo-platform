package docker

import (
	"context"
	"demo-platform/model/db"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ListImage(page int, pageSize int, order string) (*[]db.DockerImage, error) {
	var list [] db.DockerImage
	err := db.Engine.Limit(pageSize, page).OrderBy(order).Find(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}


func RemoveImage(imageID string) error{
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	_, err = cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: false,
	})
	
	return err
}