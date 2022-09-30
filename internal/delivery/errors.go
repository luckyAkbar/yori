package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrNotFound   = echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	ErrInternal   = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusInternalServerError))
)
