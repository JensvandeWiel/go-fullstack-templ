package main

import (
	"github.com/JensvandeWiel/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

// API routes are not protected by CSRF because they are used by external services.
var ignoreList = []string{"api", "swagger"}

func createRoutes(e *echo.Echo) {
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-CSRF-TOKEN",
		Skipper: func(c echo.Context) bool {
			// Skip the CSRF if the path is in the ignore list
			for _, ignorePath := range ignoreList {
				if strings.HasPrefix(c.Path(), ignorePath) {
					return true
				}
			}
			return false
		},
	}))
	iih := handlers.NewIndexInertiaHandler()
	e.GET("", iih.IndexInertiaHandle)
	{
		api := e.Group("/api")
		{
			helloWorldHandler := handlers.NewHelloWorldHandler()
			v1 := api.Group("/v1")
			v1.GET("/hello", helloWorldHandler.HelloWorldHandle)
		}
	}
}
