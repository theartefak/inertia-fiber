package inertia

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theartefak/inertia-fiber/utils"
)

// Middleware returns a middleware function that sets the version header and context for the engine.
func (e *Engine) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if e.config.AssetsPath == "" {
			panic("please provide an assets path")
		}

		// Compute the hash of the assets directory.
		hash := utils.HashDir(e.config.AssetsPath)

		// Set Vary Header to X-Inertia
		c.Set("Vary", HeaderPrefix)

		// If the request is an XHR GET request and the version header does not match the hash, return a conflict error.
		if c.Method() == "GET" && c.XHR() && c.Get(HeaderVersion, "1") != hash {
			c.Set(HeaderLocation, c.OriginalURL())
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{})
		}

		// Check if the HTTP method used in the current request is one of "PUT", "PATCH", or "DELETE"
		// and the HTTP response status code is "Found" (HTTP 302).
		exist, _ := utils.InArray(c.Method(), []string{"PUT", "PATCH", "DELETE"})
		if exist && c.Response().StatusCode() == fiber.StatusFound {
		    c.Status(fiber.StatusSeeOther)
		}

		// Check if the request header with key X-Inertia is set to the string "true".
		if c.Get(HeaderPrefix) == "true" {
		    c.Set("Vary", "Accept")
		    c.Set(HeaderPrefix, "true")
		}

		// Set the version header and context for the engine.
		c.Set(HeaderVersion, hash)
		e.ctx = c

		return c.Next()
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
