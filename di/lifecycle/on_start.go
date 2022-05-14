package lifecycle

import (
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/core/app"
	"stockcontent-monitor-demo-back/core/echo/controller"
)

func ProvidesOnStart(e *echo.Echo, binders controller.Controllers) app.OnStart {
	return func() error {
		// controller bind
		for _, b := range binders {
			b.Bind(e)
		}

		return nil
	}
}
