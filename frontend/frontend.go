package frontend

import (
	"crypto/sha256"
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
	"html/template"
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

	//go:embed public/build/.vite/manifest.json
	manifest string

	distDirFS = echo.MustSubFS(dist, "public/build")

	// Endpoints to ignore the proxy for
	ignoreList = []string{"api", "swagger"}

	inertia *gonertia.Inertia
)

const (
	manifestPath = "frontend/public/build/.vite/manifest.json"
)

func createHash() string {
	hash := sha256.New()
	hash.Write(rootTemplate)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func RegisterFrontend(e *echo.Echo) {

	var err error
	inertia, err = gonertia.NewFromBytes(rootTemplate, gonertia.WithVersion(createHash()))
	if err != nil {
		logger.GetLogger().Error("Failed to initialize Inertia", err)
		os.Exit(1)
	}
	inertia.ShareTemplateFunc("vite", vite(manifestPath, ""))
	inertia.ShareTemplateFunc("viteHead", viteHead())

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

func viteHead() func() template.HTML {
	return func() template.HTML {
		if conf, ok := config.GetConfig(); ok && conf.Environment == "dev" {
			return "<script type=\"module\" src=\"/@vite/client\"></script>"
		} else {
			return ""
		}
	}
}

func GetInertia() *gonertia.Inertia {
	return inertia
}
