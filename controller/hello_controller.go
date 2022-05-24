package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"stockcontent-monitor-demo-back/service"
)

func CreateHello(c echo.Context) error {
	var binder struct {
		Name string `json:"name"`
	}
	err := c.Bind(&binder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := service.CreateHello(binder.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": res.Id,
	})
}

func FetchHello(c echo.Context) error {
	items, err := service.FetchHello()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(items) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": items,
	})
}

func UpdateHello(c echo.Context) error {
	var binder struct {
		HelloId uuid.UUID `json:"-" param:"id"`
		Name    string    `json:"name"`
	}

	err := c.Bind(&binder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = service.UpdateHello(binder.HelloId, binder.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteHello(c echo.Context) error {
	var binder struct {
		HelloId uuid.UUID `json:"-" param:"id"`
	}
	err := c.Bind(&binder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = service.DeleteHello(binder.HelloId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
