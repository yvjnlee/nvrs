package utils

import "github.com/labstack/echo/v4"

// JSONError returns a consistent JSON error response
func JSONError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}
