package main

import (
	"demo-platform/conf"
	"demo-platform/model/db"
	"log"
)

func main() {
	r := conf.SetupRouter()
	err := db.SetupDatabase()
	if err!=nil {
		log.Fatalf("%v", err)
	}
	_ = r.Run(":8000")

}
