package main

import (
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/controller"
	"stockcontent-monitor-demo-back/env"
)

func main() {
	e := echo.New()
	e.POST("/hello", controller.CreateHello)
	e.GET("/hello", controller.FetchHello)
	e.PUT("/hello/:id", controller.UpdateHello)
	e.DELETE("/hello/:id", controller.DeleteHello)
	e.Start(env.ServeAddr)
}
