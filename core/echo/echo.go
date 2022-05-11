package echo

import (
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/core/config"
)

func New(cfg config.Config) (e *echo.Echo) {
	e = echo.New()
	return
}
