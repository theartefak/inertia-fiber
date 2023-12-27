package inertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/theartefak/inertia-fiber/utils"
)

// Render renders the Inertia component.
func (e *Engine) Render(w io.Writer, component string, props any, paths ...string) error {
	// Type assertion for props to ensure it is a fiber.Map
	propsMap, ok := props.(fiber.Map)
	if !ok {
		return fmt.Errorf("X-Inertia: props must be of type fiber.Map")
	}

	return e.partialReload(component, propsMap, w)
}

// partialReload handles the reload of the Inertia component.
func (e *Engine) partialReload(component string, props fiber.Map, w io.Writer) error {
	// Initialize an empty map for partial data
	only := make(map[string]string)

	// Retrieve the partial data from the request header
	partial := e.ctx.Get(HeaderPartialData)

	// Check if partial data exists and matches the current component
	if partial != "" && e.ctx.Get(HeaderPartialComponent) == component {
		// Populate the 'only' map with values from the partial data
		for _, value := range strings.Split(partial, ",") {
			only[value] = value
		}
	}

	// Create a new Inertia page with default values
	data := &Page{
		Component : component,
		Props     : make(fiber.Map),
		URL       : e.ctx.OriginalURL(),
		Version   : e.version,
	}

	// Copy values from the 'next' map to the Inertia page props
	for key, value := range e.next {
		if _, ok := only[key]; len(only) == 0 || ok {
			data.Props[key] = value
		}
	}

	// Copy values from the current props to the Inertia page props
	for key, value := range props {
		if _, ok := only[key]; len(only) == 0 || ok {
			data.Props[key] = value
		}
	}

	// Copy values from the context props to the Inertia page props
	contextProps := e.ctx.Context().Value(ContextKeyProps)
	if contextProps != nil {
		contextProps, ok := contextProps.(fiber.Map)
		if !ok {
			return fmt.Errorf("X-Inertia: could not convert context props to map")
		}
		for key, value := range contextProps {
			if _, ok := only[key]; len(only) == 0 || ok {
				data.Props[key] = value
			}
		}
	}

	// Reset the 'next' map for the next rendering cycle
	e.next = map[string]any{}

	// Check if XHR and Inertia should render JSON response
	renderJSON, err := strconv.ParseBool(e.ctx.Get(HeaderPrefix, "false"))
	if err != nil {
		return fmt.Errorf("X-Inertia not parsable: %w", err)
	}

	if renderJSON && e.ctx.XHR() {
		// Set headers for JSON response and return JSON representation of the Inertia page
		e.ctx.Set("Vary", "Accept")
		e.ctx.Set("X-Inertia", "true")
		e.ctx.Set("Content-Type", "application/json")
		return jsonResponse(e.ctx, data)
	}

	// Render the HTML response using the configured template and parameters
	return e.renderHTML(data, w, e.config.Template, e.Engine.Render, e.params)
}

// renderHTML prepares the data for rendering and invokes the specified renderer.
func (e *Engine) renderHTML(data *Page, w io.Writer, tmpl string, renderer func(io.Writer, string, any, ...string) error, params map[string]any) error {
	// Marshal Inertia page data to JSON
	componentData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Create a new Ziggy instance
	ziggy := utils.NewZiggy(e.ctx)

	// Marshal Ziggy data to JSON
	ziggyData, err := json.Marshal(ziggy)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Construct values for rendering the HTML page
	vals := fiber.Map{
		"Inertia" : template.HTML(fmt.Sprintf("<div id='app' data-page='%s'></div>", string(componentData))),
		"Ziggy"   : template.HTML(fmt.Sprintf("<script>const Ziggy = %s;</script>", string(ziggyData))),
		"Vite"    : utils.Vite([]string{e.config.AssetsPath + "/app.js", e.config.AssetsPath + "/Pages/" + data.Component + ".vue"}),
	}

	// Add additional parameters for rendering
	for key, value := range params {
		vals[key] = value
	}

	// Set Vary Header to X-Inertia
	e.ctx.Set("Vary", HeaderPrefix)

	// Set the Content-Type header for HTML response
	e.ctx.Set("Content-Type", "text/html")

	// Invoke the specified renderer to render the HTML page
	return renderer(w, tmpl, vals)
}

// jsonResponse sends a JSON response using Fiber for Inertia requests.
func jsonResponse(c *fiber.Ctx, page *Page) error {
	// Marshal the Inertia page to JSON
	jsonByte, err := json.MarshalIndent(page, "", "    ")
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Send the JSON response with Fiber
	return c.Status(fiber.StatusOK).JSON(string(jsonByte))
}
