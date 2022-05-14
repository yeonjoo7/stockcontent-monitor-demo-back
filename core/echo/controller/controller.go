package controller

import "github.com/labstack/echo/v4"

type Controller interface {
	Bind(e *echo.Echo)
}

type Controllers []Controller
