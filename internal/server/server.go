package server

import (
	"macro-analyst/internal/fred"
	"macro-analyst/internal/ws"

	"github.com/gofiber/fiber/v2"
)

// FiberServer wraps the Fiber application with WebSocket Hub integration.
type FiberServer struct {
	*fiber.App

	// Hub is the WebSocket message broker for real-time updates
	Hub *ws.Hub

	// FREDClient is the client for fetching macroeconomic data
	FREDClient fred.Client
}

// Config holds the configuration for the FiberServer.
type Config struct {
	ServerHeader string
	AppName      string
	FREDAPIKey   string
}

// DefaultConfig returns the default server configuration.
func DefaultConfig() Config {
	return Config{
		ServerHeader: "macro-analyst",
		AppName:      "macro-analyst",
	}
}

// New creates a new FiberServer instance with the given Hub and configuration.
func New(hub *ws.Hub, cfg ...Config) *FiberServer {
	config := DefaultConfig()
	if len(cfg) > 0 {
		config = cfg[0]
	}

	var fredClient fred.Client
	if config.FREDAPIKey != "" {
		fredClient = fred.NewClient(config.FREDAPIKey)
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: config.ServerHeader,
			AppName:      config.AppName,
		}),
		Hub:        hub,
		FREDClient: fredClient,
	}

	return server
}
