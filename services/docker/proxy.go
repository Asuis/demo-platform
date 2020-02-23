package docker

import "demo-platform/model/db"

type BindProxyForm struct {
	ContainerID int64
	ContainerPort uint16
}

func BindPort(form *BindProxyForm, user *db.User) {

}