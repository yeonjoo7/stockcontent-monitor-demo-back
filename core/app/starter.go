package app

import (
	"github.com/labstack/echo/v4"
	"stockcontent-monitor-demo-back/core/config"
	"stockcontent-monitor-demo-back/ent"
)

type OnStart func() error
type OnClose func(error)

type Start func() error

func NewStarter(
	e *echo.Echo,
	db *ent.Client,
	cfg config.Config,
	onStart OnStart,
	onClose OnClose,
) Start {
	return func() (err error) {
		err = onStart()
		if err != nil {
			return err
		}

		defer func() {
			err = db.Close()
			onClose(err)
		}()

		err = e.Start(cfg.ServerAddress())
		return
	}
}
