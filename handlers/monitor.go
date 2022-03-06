package handlers

import (
	"mon-tool-be/models"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	_ "mon-tool-be/docs"
)

type Monitor struct {
	DB *gorm.DB
}

type StoreInput struct {
	Name          string `json:"name" example:"test1"`
	Type          string `json:"type" example:"https"`
	Url           string `json:"url" example:"google.com.my"`
	RequestMethod string `json:"request_method" example:"post"`
}

// Store godoc
// @Summary      This API is to store a monitor
// @Description  If you want to create a monitor, this is the API endpoint that you should use.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authentication header"
// @Param        Body body StoreInput true "Json Body"
// @Router       /monitors [post]
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

// Edit godoc
// @Summary      This API is to edit a monitor
// @Description  If you want to edit a monitor, this is the API endpoint that you should use.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authentication header"
// @Param        monitor_id path int true "Monitor ID"
// @Param        Body body StoreInput true "Json Body"
// @Router       /monitors/{monitor_id} [patch]
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

// List godoc
// @Summary      This API is to list all monitors
// @Description  If you want to list all monitors, this is the API endpoint that you should use.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authentication header"
// @Router       /monitors [get]
func (m Monitor) List(ec echo.Context) error {
	var monitors []models.Monitor

	op := m.DB.Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).Find(&monitors)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitors)
}

// View godoc
// @Summary      This API is to view a monitor
// @Description  If you want to view a monitor, this is the API endpoint that you should use.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        monitor_id path int true "Monitor ID"
// @Param        Authorization header string true "Authentication header"
// @Router       /monitors/{monitor_id} [get]
func (m Monitor) View(ec echo.Context) error {
	id := ec.Param("id")

	monitor := models.Monitor{}

	op := m.DB.Where("id = ?", id).Where(models.Monitor{Owner: ec.Request().Header.Get("UID")}).First(&monitor)
	if op.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, op.Error.Error())
	}

	return ec.JSON(http.StatusOK, monitor)
}

// Delete godoc
// @Summary      This API is to delete a monitor
// @Description  If you want to delete a monitor, this is the API endpoint that you should use.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        monitor_id path int true "Monitor ID"
// @Param        Authorization header string true "Authentication header"
// @Router       /monitors/{monitor_id} [delete]
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
