package handlers

import (
	"mon-tool-be/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Record struct {
	DB *gorm.DB
}

func (m Record) View(ec echo.Context) error {
	id := ec.Param("id")

	monitor := models.Monitor{}

	op := m.DB.Where("id = ?", id).Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).Preload("Records").First(&monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor.Records)
}
