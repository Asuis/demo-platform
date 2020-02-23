package db

type Display struct {
	ID int64
	ContainerID int64
	ExposeUrl int64
	Group *[] DisplayGroup
}

type DisplayGroup struct {
	Name string
	Containers *[]DockerContainer
	Displays *[] Display
	OwnerID int64
	Owner User
}