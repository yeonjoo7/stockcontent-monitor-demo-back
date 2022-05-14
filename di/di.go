package di

import (
	"github.com/google/wire"
	"stockcontent-monitor-demo-back/core/app"
	"stockcontent-monitor-demo-back/core/config"
	"stockcontent-monitor-demo-back/core/db"
	"stockcontent-monitor-demo-back/core/echo"
	"stockcontent-monitor-demo-back/core/echo/controller"
	lifecycle2 "stockcontent-monitor-demo-back/di/lifecycle"
	"stockcontent-monitor-demo-back/di/value"
	"stockcontent-monitor-demo-back/hello/handler"
	"stockcontent-monitor-demo-back/hello/repository"
	"stockcontent-monitor-demo-back/hello/usecase"
)

var DI = wire.NewSet(
	app.NewStarter,
	lifecycle2.ProvidesOnStart,
	lifecycle2.ProvidesOnClose,
	CastControllers,
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

// repository
var repositoryProviders = wire.NewSet(
	repository.NewHelloMySQLRepository,
)

// useCase
var useCaseProviders = wire.NewSet(
	usecase.NewHelloUseCase,
)

func CastControllers(
	c1 *handler.HelloController,
) controller.Controllers {
	return controller.Controllers{
		c1,
	}
}

// controller
var controllerProviders = wire.NewSet(
	handler.HelloControllerProvider,
)

var staticValue = wire.NewSet(
	wire.InterfaceValue(new(config.Config), config.Default),
	value.ProvidesTimeout,
)
