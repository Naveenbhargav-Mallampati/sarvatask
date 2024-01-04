package api

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

func Driver(logger hclog.Logger) {
	e := echo.New()
	e.Use(LoggerMiddleware(logger))
	e.POST("/upload", Upload)
	logger.Info("Started Upload Service")
	e.Start(":8000")
}

func LoggerMiddleware(logger hclog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			logger.Info("Received request", "method", c.Request().Method, "path", c.Request().URL.Path)
			return next(c)
		}
	}
}
