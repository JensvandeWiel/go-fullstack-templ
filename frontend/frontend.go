package frontend

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/JensvandeWiel/config"
	"github.com/JensvandeWiel/logger"
	_ "github.com/JensvandeWiel/logger"
	middleware2 "github.com/JensvandeWiel/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/romsar/gonertia"
	"net/url"
	"os"
	"path"
	"strings"
)

//Credit to Dan Hawkins for the below code
//https://github.com/danhawkins/go-vite-react-example/blob/main/frontend/frontend.go

var (
	//go:embed public/build/*
	dist embed.FS

	//go:embed root.gohtml
	rootTemplate []byte

	distDirFS = echo.MustSubFS(dist, "dist")

	// Endpoints to ignore the proxy for
	ignoreList = []string{"api", "swagger"}

	inertia *gonertia.Inertia
)

const (
	manifestPath = "frontend/public/build/manifest.json"
)

func RegisterFrontend(e *echo.Echo) {
	var err error
	inertia, err = gonertia.NewFromBytes(rootTemplate, gonertia.WithVersionFromFile(manifestPath))
	if err != nil {
		logger.GetLogger().Error("Failed to initialize Inertia", err)
		os.Exit(1)
	}
	inertia.ShareTemplateFunc("vite", vite(manifestPath, "dist"))

	e.Use(echo.WrapMiddleware(inertia.Middleware), middleware2.AttachInertia(inertia))

	if conf, ok := config.GetConfig(); ok && conf.Environment == "dev" {
		logger.GetLogger().Info("Running in dev mode")
		setupDevProxy(e)
		return
	}
	e.StaticFS("/", distDirFS)
	return
}

func setupDevProxy(e *echo.Echo) {
	url, err := url.Parse("http://localhost:5173/")
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

	e.Group("").Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
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

func vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		logger.GetLogger().Error("Failed to open manifest file", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)

	if err != nil {
		logger.GetLogger().Error("Failed to unmarshal vite manifest file to json", err)
	}

	return func(p string) (string, error) {

		// If in dev mode and the asset is in the viteAssets map, return the vite asset path
		if conf, ok := config.GetConfig(); ok && conf.Environment == "dev" {
			if _, ok := viteAssets[p]; ok {
				return path.Join("/", p), nil
			}
		}
		// If in prod mode and the asset is in the viteAssets map, return the dist asset path
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}

func GetInertia() *gonertia.Inertia {
	return inertia
}
