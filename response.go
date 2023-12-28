package inertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/theartefak/inertia-fiber/utils"
)

// View function returns a Fiber handler function for rendering Inertia.js views.
func (e *Engine) View(component string, props fiber.Map) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the request header indicates JSON rendering
		renderJSON, err := strconv.ParseBool(c.Get(HeaderPrefix, "false"))
		if err != nil {
			return fmt.Errorf("X-Inertia not parsable: %w", err)
		}

		// Attempt to perform a partial reload of the component
		partial, err := e.partialReload(component, props, c)
		if err != nil {
			return fmt.Errorf("X-Inertia: %w", err)
		}

		// If JSON rendering is requested and it's an XHR request, return JSON response
		if renderJSON && c.XHR() {
			c.Set("Vary", "Accept")
			c.Set("X-Inertia", "true")
			c.Set("Content-Type", "application/json")
			return c.Status(fiber.StatusOK).JSON(partial)
		}

		// Render the HTML page using the provided template and parameters
		return e.renderHTML(partial, c, e.config.Template, e.params)
	}
}

// partialReload function performs a partial reload of the component data.
func (e *Engine) partialReload(component string, props fiber.Map, c *fiber.Ctx) (*Page, error) {
	only := make(map[string]string)

	// Extract partial data from the request header
	partial := c.Get(HeaderPartialData)

	// Process partial data if present and matches the current component
	if partial != "" && c.Get(HeaderPartialComponent) == component {
		for _, value := range strings.Split(partial, ",") {
			only[value] = value
		}
	}

	// Create a new Page object with relevant data for the component
	data := &Page{
		Component: component,
		Props:     make(fiber.Map),
		URL:       c.OriginalURL(),
		Version:   e.version,
	}

	// Copy the next data for keys specified in the partial reload
	for key, value := range e.next {
		if _, ok := only[key]; len(only) == 0 || ok {
			data.Props[key] = value
		}
	}

	// Copy the props for keys specified in the partial reload
	for key, value := range props {
		if _, ok := only[key]; len(only) == 0 || ok {
			data.Props[key] = value
		}
	}

	// Copy the context props for keys specified in the partial reload
	contextProps := c.Context().Value("props")
	if contextProps != nil {
		contextProps, ok := contextProps.(fiber.Map)
		if !ok {
			return nil, fmt.Errorf("X-Inertia: could not convert context props to map")
		}
		for key, value := range contextProps {
			if _, ok := only[key]; len(only) == 0 || ok {
				data.Props[key] = value
			}
		}
	}

	// Reset the next data to an empty map
	e.next = map[string]any{}

	return data, nil
}

// renderHTML function renders the HTML page using the specified template and parameters.
func (e *Engine) renderHTML(data *Page, c *fiber.Ctx, tmpl string, params map[string]any) error {
	// Marshal the component data to JSON
	componentData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Create a Ziggy instance for client-side routing
	ziggy := utils.NewZiggy(c)

	// Marshal Ziggy data to JSON
	ziggyData, err := json.Marshal(ziggy)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Define values for the template rendering
	vals := fiber.Map{
		"Inertia": template.HTML(fmt.Sprintf("<div id='app' data-page='%s'></div>", string(componentData))),
		"Ziggy":   template.HTML(fmt.Sprintf("<script>const Ziggy = %s;</script>", string(ziggyData))),
		"Vite":    utils.Vite([]string{e.config.AssetsPath + "/app.js", e.config.AssetsPath + "/Pages/" + data.Component + ".vue"}),
	}

	// Include additional parameters for template rendering
	for key, value := range params {
		vals[key] = value
	}

	// Set response headers
	c.Set("Vary", HeaderPrefix)
	c.Set("Content-Type", "text/html")

	// Render the HTML page using the specified template and values
	return e.Render(c, tmpl, vals)
}
