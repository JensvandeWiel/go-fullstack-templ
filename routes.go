package main

import (
	"github.com/JensvandeWiel/handlers"
	"github.com/labstack/echo/v4"
)

func createRoutes(e *echo.Echo) {
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
