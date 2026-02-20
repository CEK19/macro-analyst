package server

import (
	"context"
	"time"

	"macro-analyst/internal/fred"

	"github.com/gofiber/fiber/v2"
)

const (
	// RequestTimeout for FRED API calls.
	RequestTimeout = 10 * time.Second
)

// GetAllTickersHandler returns all available FRED tickers with descriptions.
func (s *FiberServer) GetAllTickersHandler(c *fiber.Ctx) error {
	if s.FREDClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "FRED API client not configured",
		})
	}

	tickers := fred.AllTickers()
	response := make([]fiber.Map, len(tickers))

	for i, ticker := range tickers {
		response[i] = fiber.Map{
			"symbol":      ticker.String(),
			"description": ticker.Description(),
		}
	}

	return c.JSON(fiber.Map{
		"tickers": response,
		"count":   len(response),
	})
}

// GetTickerDataHandler returns historical observations for a specific ticker.
func (s *FiberServer) GetTickerDataHandler(c *fiber.Ctx) error {
	if s.FREDClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "FRED API client not configured",
		})
	}

	symbol := c.Params("symbol")
	ticker := fred.Ticker(symbol)

	// Parse query parameters
	opts := &fred.QueryOptions{
		StartDate: c.Query("start_date", ""),
		EndDate:   c.Query("end_date", ""),
		Limit:     c.QueryInt("limit", fred.DefaultLimit),
		SortOrder: c.Query("sort_order", "desc"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	data, err := s.FREDClient.GetSeriesObservations(ctx, ticker, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

// GetLatestValueHandler returns the most recent value for a specific ticker.
func (s *FiberServer) GetLatestValueHandler(c *fiber.Ctx) error {
	if s.FREDClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "FRED API client not configured",
		})
	}

	symbol := c.Params("symbol")
	ticker := fred.Ticker(symbol)

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	latest, err := s.FREDClient.GetLatestValue(ctx, ticker)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(latest)
}

// GetAllLatestHandler returns the latest values for all supported tickers.
func (s *FiberServer) GetAllLatestHandler(c *fiber.Ctx) error {
	if s.FREDClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "FRED API client not configured",
		})
	}

	tickers := fred.AllTickers()

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	result, err := s.FREDClient.GetMultipleLatest(ctx, tickers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
