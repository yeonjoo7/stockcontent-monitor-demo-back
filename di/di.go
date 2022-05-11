package di

import (
	"github.com/google/wire"
	"stockcontent-monitor-demo-back/core/app"
	"stockcontent-monitor-demo-back/core/config"
	"stockcontent-monitor-demo-back/core/db"
	"stockcontent-monitor-demo-back/core/echo"
	"stockcontent-monitor-demo-back/core/echo/binder"
	"stockcontent-monitor-demo-back/core/lifecycle"
	"stockcontent-monitor-demo-back/di/provides"
	"stockcontent-monitor-demo-back/hello/handler"
	"stockcontent-monitor-demo-back/hello/repository"
	"stockcontent-monitor-demo-back/hello/usecase"
)

var DI = wire.NewSet(
	app.NewStarter,
	lifecycle.ProvidesOnStart,
	lifecycle.ProvidesOnClose,
	binder.ProvidesController,
	staticValue,
	repositoryProviders,
	useCaseProviders,
	controllerProviders,
	infraProviders,
)

var infraProviders = wire.NewSet(
	echo.New,
	db.New,
)

var repositoryProviders = wire.NewSet(
	repository.NewHelloMySQLRepository,
)

var useCaseProviders = wire.NewSet(
	usecase.NewHelloUseCase,
)

var controllerProviders = wire.NewSet(
	handler.HelloControllerProvider,
)

var staticValue = wire.NewSet(
	wire.InterfaceValue(new(config.Config), config.Default),
	provides.ProvidesTimeout,
)
