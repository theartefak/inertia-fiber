package inertia

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type (
	// Config represents the configuration for the Inertia engine.
	Config struct {
		Root       string          // The root directory of the application.
		FS         http.FileSystem // The file system to use for loading templates and assets.
		AssetsPath string          // The path to the assets directory.
		Template   string          // The name of the template to use.
	}

	// Engine represents the Inertia engine.
	Engine struct {
		*html.Engine
		ctx    *fiber.Ctx     // The current context.
		config Config         // The configuration.
		props  map[string]any // The current props.
		next   map[string]any // The next props.
		params map[string]any // The current params.
	}
)

// DefaultConfig is the default config
var DefaultConfig = Config{
	Root       : "resources/views",
	AssetsPath : "resources/js",
	Template   : "app",
}
