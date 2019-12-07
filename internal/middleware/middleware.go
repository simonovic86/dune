package middleware

import "github.com/labstack/echo"

type GoMiddleware struct {
	// may be needed by middleware
}

// CORS interceptor
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// initialize middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}