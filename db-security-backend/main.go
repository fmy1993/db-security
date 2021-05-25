package main

import (
	"db-security-backend/config"
	"log"

	"db-security-backend/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.ParseConfig("./config/engine.json")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	_, err = config.OrmEngine(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	engine := gin.Default()
	RegisterRouter(engine)
	_ = engine.Run(cfg.AppHost + ":" + cfg.AppPort)
}

func RegisterRouter(engine *gin.Engine) {
	new(controller.UserController).Router(engine)
	new(controller.StaffController).Router(engine)
	new(controller.IpController).Router(engine)
	new(controller.DownloadRecordController).Router(engine)
}
