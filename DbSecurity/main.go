package main

import (
	"log"

	"DbSecurity/controller"
	"DbSecurity/middleware"
	"DbSecurity/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := tool.ParseConfig("./config/engine.json")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	_, err = tool.OrmEngine(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	engine := gin.Default()
	engine.Use(middleware.Cors())
	RegisterRouter(engine)
	_ = engine.Run(cfg.AppHost + ":" + cfg.AppPort)
}

func RegisterRouter(engine *gin.Engine) {
	new(controller.UserController).Router(engine)
	new(controller.PatientController).Router(engine)
	new(controller.IpController).Router(engine)
	new(controller.ImplantIndexController).Router(engine)
}
