package handlers

import (
	"mon-tool-be/models"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

	// add UID as owner
	monitor.Owner = ec.Request().Header.Get("UID")

	// fix domain
	u, err := url.Parse(monitor.Url)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	monitor.Url = u.Host + u.Path + u.RawQuery + u.Fragment

	op := m.DB.Create(monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor)
}

func (m Monitor) Edit(ec echo.Context) error {
	monitor := new(models.Monitor)

	id := ec.Param("id")

	op := m.DB.Where("id = ?", id).Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).First(&monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	err := ec.Bind(monitor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "binding error")
	}

	err = monitor.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// add UID as owner
	monitor.Owner = ec.Request().Header.Get("UID")

	// fix domain
	u, err := url.Parse(monitor.Url)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	monitor.Url = u.Host + u.Path + u.RawQuery + u.Fragment

	op = m.DB.Save(monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor)
}

func (m Monitor) List(ec echo.Context) error {
	var monitors []models.Monitor

	op := m.DB.Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).Find(&monitors)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitors)
}

func (m Monitor) View(ec echo.Context) error {
	id := ec.Param("id")

	monitor := models.Monitor{}

	op := m.DB.Where("id = ?", id).Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).First(&monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor)
}

func (m Monitor) Delete(ec echo.Context) error {
	id := ec.Param("id")

	monitor := models.Monitor{}

	op := m.DB.Where("id = ?", id).Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).First(&monitor).Delete(&monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, map[string]interface{}{
		"message": "monitor deleted",
	})
}
