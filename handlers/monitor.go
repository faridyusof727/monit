package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"mon-tool-be/models"
	"net/http"
)

type Monitor struct {
	DB *gorm.DB
}

func (m Monitor) Store(ec echo.Context) error {
	monitor := new(models.Monitor)

	err := ec.Bind(monitor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = monitor.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	op := m.DB.Create(monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor)
}
