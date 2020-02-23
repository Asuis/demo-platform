package main

import (
	"demo-platform/conf"
	"demo-platform/model/db"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := conf.SetupRouter()
	err := db.SetupDatabase()
	if err!=nil {
		log.Fatalf("%v", err)
	}
	gin.SetMode("debug")
	_ = r.Run(":8000")

}
