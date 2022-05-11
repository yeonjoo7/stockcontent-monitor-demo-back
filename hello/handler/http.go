package handler

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"net/http"
	"stockcontent-monitor-demo-back/domain"
)

const (
	tag = "HELLO-CONTROLLER"
)

var HelloControllerProvider = wire.NewSet(
	wire.Struct(new(HelloController), "*"),
)

type HelloController struct {
	UseCase domain.HelloUseCase
}

func (c *HelloController) Bind(e *echo.Echo) {
	e.GET("/", c.hello)
}

func (c *HelloController) hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "hello world",
	})
}
