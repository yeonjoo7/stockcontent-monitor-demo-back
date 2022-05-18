package lifecycle

import (
	"context"
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/core/app"
	"stockcontent-monitor-demo-back/core/config"
	"stockcontent-monitor-demo-back/core/echo/controller"
	"stockcontent-monitor-demo-back/ent"
	"stockcontent-monitor-demo-back/ent/migrate"
)

func ProvidesOnStart(
	cfg config.Config,
	e *echo.Echo,
	controllers controller.Controllers,
	db *ent.Client,
) app.OnStart {
	return func() (err error) {
		// controller bind
		for _, b := range controllers {
			b.Bind(e)
		}

		if cfg.IsDebug() {
			err = db.Schema.Create(context.Background(),
				migrate.WithDropIndex(true),
				migrate.WithDropColumn(true),
				migrate.WithForeignKeys(true),
			)

			return
		}

		return
	}
}
