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
	// Merge the current props with the next props.
	dt := e.props
	for k, v := range e.next {
		dt[k] = v
	}

	// Merge the current props with the specified props.
	for k, v := range props {
		dt[k] = v
	}

	// Create a new data map with the component, props, URL, and version.
	data := map[string]interface{}{
		"component" : component,
		"props"     : dt,
		"url"       : e.ctx.OriginalURL(),
		"version"   : e.ctx.Get(HeaderVersion, ""),
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
		return jsonResponse(e.ctx, data)
	}

	// Otherwise, return an HTML response.
	return e.toResponse(data, w, e.config.Template, e.Engine.Render, e.params)
}

// toResponse returns an HTML response with the specified data and template.
func (e *Engine) toResponse(data fiber.Map, w io.Writer, tmpl string, renderer func(io.Writer, string, any, ...string) error, params map[string]any) error {
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
		"Vite"    : utils.Vite([]string{e.config.AssetsPath + "/app.js", e.config.AssetsPath + "/Pages/" + data["component"].(string) + ".vue"}),
	}

	// Merge the params into the map.
	for k, v := range params {
		vals[k] = v
	}

	// Render the template with the map.
	return renderer(w, tmpl, vals)
}

// jsonResponse returns a JSON response with the specified page.
func jsonResponse(c *fiber.Ctx, page fiber.Map) error {
	// Marshal the page to JSON.
	jsonByte, _ := json.Marshal(page)

	// Return a JSON response with the page.
	return c.Status(fiber.StatusOK).JSON(string(jsonByte))
}

// partialReload reloads the component with the specified props.
func partialReload(c *fiber.Ctx, component string, props fiber.Map) fiber.Map {
	// If the component is the same as the partial component, create a new map with the partial props.
	if c.Get(HeaderPartialComponent, "/") == component {
		var newProps = make(fiber.Map)
		partials := strings.Split(c.Get(HeaderPartialData, ""), ",")

		for _, partial := range partials {
			if val, ok := props[partial]; ok {
				newProps[partial] = val
			}
		}

		if len(newProps) > 0 {
			return newProps
		}
	}

	// Return the original props.
	return props
}
