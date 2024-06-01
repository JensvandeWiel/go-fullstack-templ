package frontend

import (
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-fullstack-templ/config"
	"go-fullstack-templ/logger"
	_ "go-fullstack-templ/logger"
	"net/url"
	"os"
	"strings"
)

//Credit to Dan Hawkins for the below code
//https://github.com/danhawkins/go-vite-react-example/blob/main/frontend/frontend.go

var (
	//go:embed dist/*
	dist embed.FS

	//go:embed dist/index.html
	indexHTML embed.FS

	distDirFS     = echo.MustSubFS(dist, "dist")
	distIndexHTML = echo.MustSubFS(indexHTML, "dist")

	// Endpoints to ignore the proxy for
	ignoreList = []string{"/api", "/swagger"}
)

func RegisterFrontend(e *echo.Echo) {
	if env, ok := config.GetKey("evnironment"); ok && env.(string) == "dev" {
		logger.GetLogger().Info("Running in dev mode")
		setupDevProxy(e)
		return
	}
	// Use the static assets from the dist directory
	e.FileFS("/", "index.html", distIndexHTML)
	e.StaticFS("/", distDirFS)
}

func setupDevProxy(e *echo.Echo) {
	url, err := url.Parse("http://localhost:5173")
	if err != nil {
		logger.GetLogger().Error("Failed to parse the URL for the dev server", err, url)
		os.Exit(1)
	}
	// Setup a proxy to the vite dev server on localhost:5173
	balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: url,
		},
	})

	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: balancer,
		Skipper: func(c echo.Context) bool {
			// Skip the proxy if the path is in the ignore list
			for _, ignorePath := range ignoreList {
				if strings.HasPrefix(c.Path(), ignorePath) {
					return true
				}
			}
			return false
		},
	}))
}
