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

		// If the request is an XHR GET request and the version header does not match the hash, return a conflict error.
		if c.Method() == "GET" && c.XHR() && c.Get(HeaderVersion, "1") != hash {
			c.Set(HeaderLocation, c.Path())
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{})
		}

		// Set the version header and context for the engine.
		c.Set(HeaderVersion, hash)
		e.ctx = c

		return c.Next()
	}
}
