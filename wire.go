//go:build wireinject

package main

import (
	"github.com/google/wire"
	"stockcontent-monitor-demo-back/core/app"
	"stockcontent-monitor-demo-back/di"
)

func getStarter() app.Start {
	wire.Build(di.DI)
	return nil
}
