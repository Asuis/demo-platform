package group


type Group struct {
	ID int64
	GroupName string
	Type int64
	RepoID int64
	ContainerID int64
}

func CreateGroup() {}

func JoinGroup() {}

func DeleteGroup() {}