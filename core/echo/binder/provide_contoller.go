package binder

import (
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/hello/handler"
)

type Binder interface {
	Bind(e *echo.Echo)
}

type Binders []Binder

func ProvidesController(
	controller *handler.HelloController,
) Binders {
	return castBinders(controller)
}

func castBinders(controllers ...Binder) Binders {
	return controllers
}
