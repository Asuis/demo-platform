package display

import "demo-platform/model/db"

type Queue struct {
	Name string
	ContainerIDs [] int64
}

func CreateDisplayQueue(form Queue) {
	//检查 权限

	return
}

type QueueActionForm struct {
	Name string
	NextID int64
}
func NextDisplay(form *QueueActionForm) {
	//依次启动nextID 前后3个容器配置
	var list [] db.DisplayGroup
	db.Engine.Where("ContainerID = ? and Di", form.NextID).Limit(3).OrderBy("ID").Find(&list)
}

func GetDisplayStatus() {

}
