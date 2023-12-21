package utils

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Ziggy is a struct to export compatible Echo routes for https://github.com/tightenco/ziggy.
type (
	Ziggy struct {
		BaseDomain   string                `json:"domain"`
		BasePort     int                   `json:"port"`
		BaseProtocol string                `json:"protocol"`
		BaseUrl      string                `json:"url"`
		Group        string                `json:"group"`
		Routes       map[string]ZiggyRoute `json:"routes"`
	}

// ZiggyRoute represents a single route for https://github.com/tightenco/ziggy.
	ZiggyRoute struct {
		Uri     string   `json:"uri"`
		Methods []string `json:"methods"`
		Domain  string   `json:"domain"`
	}
)

// NewZiggy creates a new Ziggy instance.
func NewZiggy(c *fiber.Ctx) Ziggy {
	var z Ziggy

	// Set the base protocol to the protocol used by the request.
	z.BaseProtocol = c.Protocol()

	// Get the hostname from the request.
	host := c.Hostname()

	// Split the hostname into its components.
	splitHost := strings.Split(host, ":")

	// Set the base domain to the first component of the hostname.
	z.BaseDomain = splitHost[0]

	// Set the base URL to the protocol and domain.
	z.BaseUrl = z.BaseProtocol + "://" + z.BaseDomain

	// If the hostname contains a port number, set the base port and add it to the base URL.
	if len(splitHost) > 1 {
		port, err := strconv.Atoi(splitHost[1])
		if err == nil && port > 0 {
			z.BasePort = port
			z.BaseUrl += ":" + strconv.Itoa(z.BasePort)
		}
	}

	// Create a new map to store the routes.
	z.Routes = make(map[string]ZiggyRoute, len(c.App().GetRoutes()))

	// Iterate over the routes and add them to the map.
	for _, route := range c.App().GetRoutes() {
		if route.Name == "" && route.Path == "/" {
			continue
		}

		if ziggyRoute, ok := z.Routes[route.Name]; ok {
			if !contains(ziggyRoute.Methods, route.Method) {
				ziggyRoute.Methods = append(ziggyRoute.Methods, route.Method)
			}
			z.Routes[route.Name] = ziggyRoute
		} else {
			z.Routes[route.Name] = ZiggyRoute{
				Uri:     route.Path,
				Methods: []string{route.Method},
			}
		}
	}

	return z
}

// contains returns true if the slice contains the element.
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
