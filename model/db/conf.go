package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"xorm.io/core"
	"xorm.io/xorm"
)

type DatabaseConf struct {
	host string
	user string
	pass string
}
var (
	Engine *xorm.Engine
	tables [] interface{}
)

func init() {
	tables = append(tables,
		new(User),new(EmailAddress),
		new(Repository),new(Branch),new(Access),new(Collaboration),new(Mirror),new(DockerImage), new(DockerContainer))

	gonicNames := []string{"SSL"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}
}

func SetupDatabase() error {
	var err error
	Engine, err = xorm.NewEngine("mysql", "root:Mysql127117@/beta?charset=utf8")

	if err != nil {
		log.Fatalf("")
	}
	//setting  fmapper
	Engine.SetMapper(core.SameMapper{})
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "beta_")
	Engine.SetTableMapper(tbMapper)
	if err = Engine.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync structs to database tables: %v\n", err)
	}
	return nil
}

func SetEngine() (err error) {

	return nil
}
