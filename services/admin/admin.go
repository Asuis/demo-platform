package admin

import "demo-plaform/model/db"

func ListUser(page int, pageSize int, order string) (*[] db.User, error) {
	var list [] db.User
	err := db.Engine.Limit(pageSize, page).OrderBy(order).Find(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func ListContainer(page int, pageSize int, order string) (*[] db.DockerContainer, error){
	var list [] db.DockerContainer
	err := db.Engine.Limit(pageSize, page).OrderBy(order).Find(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func ListRepository(page int, pageSize int, order string) (*[] db.Repository, error) {
	var list [] db.Repository
	err := db.Engine.Limit(pageSize, page).OrderBy(order).Find(&list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

