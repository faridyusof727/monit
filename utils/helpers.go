package utils

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func JsonMapper(input io.Reader) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(input).Decode(&jsonMap)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "error parsing json structure")
	}
	return jsonMap, nil
}
