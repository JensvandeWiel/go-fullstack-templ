package main

import (
	"github.com/JensvandeWiel/handlers"
	"github.com/labstack/echo/v4"
)

func createRoutes(root *echo.Group) {
	//V1
	v1 := root.Group("/v1")
	indexHandler := handlers.NewIndexHandler()
	v1.GET("/hello", indexHandler.HelloWorldHandle)
}
