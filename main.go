package main

import (
	"fmt"
	"github.com/JensvandeWiel/config"
	"github.com/JensvandeWiel/docs"
	"github.com/JensvandeWiel/frontend"
	"github.com/JensvandeWiel/logger"
	_ "github.com/JensvandeWiel/logger"
	"github.com/JensvandeWiel/middleware"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			go-fullstack-templ
//	@version		1.0
//	@description	Fullstack golang template w ith js framework of choice

//	@contact.name	Jens van de Wiel
//	@contact.url	https://jens.vandewiel.eu
//	@contact.email	jens.vdwiel@gmail.com

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @license.name	MIT

func main() {
	conf, ok := config.GetConfig()
	if !ok {
		panic("Failed to get config")
	}
	e := createEcho()
	createRoutes(e)

	url := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	docs.SwaggerInfo.Host = url

	if conf.Environment != "production" {
		logger.GetLogger().Debug("Environment is not production, enabling swagger endpoint")
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	logger.GetLogger().Info(fmt.Sprintf("Starting server on: %s", url))
	logger.GetLogger().Error(e.Start(url).Error())
}

func createEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(slogecho.New(logger.GetLogger()), middleware.AttachRequestID())
	createRoutes(e)
	frontend.RegisterFrontend(e)
	return e
}
