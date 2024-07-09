package middleware

import (
	"github.com/JensvandeWiel/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/romsar/gonertia"
	"log/slog"
)

func AttachRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := uuid.NewRandom()
			if err != nil {
				slog.Error("Failed to generate request ID", slog.String("error", err.Error()))
				return err
			}

			c.Set("request_id", id.String())
			c.Set("logger", logger.GetLogger().With("request_id", id.String()))
			c.Response().Header().Set(echo.HeaderXRequestID, id.String())

			return next(c)
		}
	}
}

func GetRequestID(c echo.Context) string {
	return c.Get("request_id").(string)
}

func GetLogger(c echo.Context) *slog.Logger {
	return c.Get("logger").(*slog.Logger)
}

func GetInertia(c echo.Context) *gonertia.Inertia {
	return c.Get("inertia").(*gonertia.Inertia)
}

func AttachInertia(i *gonertia.Inertia) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("inertia", i)
			return next(c)
		}
	}
}
