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

// Render renders the component with the specified props.
func (e *Engine) Render(w io.Writer, component string, props any, paths ...string) error {
	// Reload the component with the specified props.
	p := partialReload(e.ctx, component, props.(fiber.Map))

	// Display the component with the specified props.
	return e.display(component, p, w)
}

// display displays the component with the specified props.
func (e *Engine) display(component string, props fiber.Map, w io.Writer) error {
	data := &Page{
		Component: component,
		Props: make(fiber.Map),
		URL: e.ctx.OriginalURL(),
		Version: e.version,
	}

	// Merge the current props with the next props.
	for key, value := range e.next {
		data.Props[key] = value
	}

	// Merge the current props with the specified props.
	for key, value := range props {
		data.Props[key] = value
	}

	// Retrieve the value associated with the ContextKeyProps key from the context.
	contextProps := e.ctx.Context().Value(ContextKeyProps)

	// Check if the retrieved contextProps is not nil.
	if contextProps != nil {
		// Attempt to type assert contextProps to a fiber.Map.
		contextProps, ok := contextProps.(fiber.Map)

		// Check if the type assertion was successful.
		if !ok {
			// If the type assertion fails, return an error indicating the conversion failure.
			return fmt.Errorf("X-Inertia: could not convert context props to map")
		}

		// Iterate over key-value pairs in the contextProps map.
		for key, value := range contextProps {
			// Copy each key-value pair from contextProps to the data.Props map.
			data.Props[key] = value
		}
	}


	// Clear the next props.
	e.next = map[string]any{}

	// Check if the response should be rendered as JSON.
	renderJSON, err := strconv.ParseBool(e.ctx.Get(HeaderPrefix, "false"))
	if err != nil {
		return fmt.Errorf("X-Inertia not parsable: %w", err)
	}

	// If the response should be rendered as JSON, return a JSON response.
	if renderJSON && e.ctx.XHR() {
		e.ctx.Set("Vary", "Accept")
		e.ctx.Set("X-Inertia", "true")
		e.ctx.Set("Content-Type", "application/json")
		return jsonResponse(e.ctx, data)
	}

	// Otherwise, return an HTML response.
	return e.toResponse(data, w, e.config.Template, e.Engine.Render, e.params)
}

// toResponse returns an HTML response with the specified data and template.
func (e *Engine) toResponse(data *Page, w io.Writer, tmpl string, renderer func(io.Writer, string, any, ...string) error, params map[string]any) error {
	// Marshal the data to JSON.
	componentData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	// Create a new Ziggy instance.
	ziggy := utils.NewZiggy(e.ctx)

	// Marshal the Ziggy instance to JSON.
	ziggyData, _ := json.Marshal(ziggy)

	// Create a new map with the Inertia key, the component data, and the Ziggy data.
	vals := fiber.Map{
		"Inertia" : template.HTML(fmt.Sprintf("<div id='app' data-page='%s'></div>", string(componentData))),
		"Ziggy"   : template.HTML(fmt.Sprintf("<script>const Ziggy = %s;</script>", string(ziggyData))),
		"Vite"    : utils.Vite([]string{e.config.AssetsPath + "/app.js", e.config.AssetsPath + "/Pages/" + data.Component + ".vue"}),
	}

	// Merge the params into the map.
	for key, value := range params {
		vals[key] = value
	}

	e.ctx.Set("Content-Type", "text/html")

	// Render the template with the map.
	return renderer(w, tmpl, vals)
}

// jsonResponse returns a JSON response with the specified page.
func jsonResponse(c *fiber.Ctx, page *Page) error {
	// Marshal the page to JSON.
	jsonByte, _ := json.Marshal(page)

	// Return a JSON response with the page.
	return c.Status(fiber.StatusOK).JSON(string(jsonByte))
}

// partialReload reloads the component with the specified props.
func partialReload(c *fiber.Ctx, component string, props fiber.Map) fiber.Map {
	// Initialize a map to store only the values needed for partial reload.
	only := make(fiber.Map)

	// Retrieve the value of the "X-Inertia-Partial-Data" header from the context.
	partial := c.Get(HeaderPartialData)

	// Check if "X-Inertia-Partial-Data" is not empty and the "X-Inertia-Partial-Component" header
	// has a value equal to the provided component name.
	if partial != "" && c.Get(HeaderPartialComponent) == component {
		// Iterate over the values obtained by splitting the "partial" string using commas.
		for _, value := range strings.Split(partial, ",") {
			// Add each value to the "only" map.
			only[value] = value
		}

		// If there are values in the "only" map, indicating a partial reload is needed,
		// return the map with only the specified values.
		if len(only) > 0 {
			return only
		}
	}

	// If no partial reload is needed, return the original properties map.
	return props
}
