package inertia

import "github.com/gofiber/template/html/v2"

// The Base "X-Inertia" header prefixes
const (
	HeaderPrefix           = "X-Inertia"
	HeaderVersion          = HeaderPrefix + "-Version"
	HeaderLocation         = HeaderPrefix + "-Location"
	HeaderPartialData      = HeaderPrefix + "-Partial-Data"
	HeaderPartialComponent = HeaderPrefix + "-Partial-Component"
)

// New creates a new instance of the Engine.
func New(config ...Config) *Engine {
	var engine *html.Engine
	cfg := DefaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]
	}

	// If no assets path is specified, use the default assets path.
	if cfg.AssetsPath == "" {
		cfg.AssetsPath = "resources/js"
	}

	// If no root path is specified, use the default root path.
	if cfg.Root == "" {
		cfg.Root = "resources/views"
	}

	// If a file system is provided, use it to create the engine.
	if cfg.FS != nil {
		engine = html.NewFileSystem(cfg.FS, ".html")
	} else {
		// Otherwise, use the root directory to create the engine.
		engine = html.New(cfg.Root, ".html")
	}

	// If no template is specified, use the default template.
	if cfg.Template == "" {
		cfg.Template = "app"
	}

	// Create a new instance of the Engine with the specified configuration.
	return &Engine{
		Engine : engine,
		config : cfg,
		props  : make(map[string]any),
		params : make(map[string]any),
		next   : make(map[string]any),
	}
}

func (e *Engine) Share(name string, value any) {
	e.props[name] = value
}

func (e *Engine) AddProp(name string, value any) {
	e.next[name] = value
}

func (e *Engine) AddParam(name string, value any) {
	e.params[name] = value
}
