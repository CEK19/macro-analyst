package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL is the FRED API base endpoint.
	BaseURL = "https://api.stlouisfed.org/fred"

	// DefaultTimeout for HTTP requests.
	DefaultTimeout = 10 * time.Second

	// DefaultLimit for observations.
	DefaultLimit = 100
)

// Client defines the interface for FRED API operations.
// This interface allows for easy mocking in tests.
type Client interface {
	GetSeriesObservations(ctx context.Context, ticker Ticker, opts *QueryOptions) (*SeriesData, error)
	GetLatestValue(ctx context.Context, ticker Ticker) (*LatestValue, error)
	GetMultipleLatest(ctx context.Context, tickers []Ticker) (*MultiTickerResponse, error)
	GetSeriesInfo(ctx context.Context, ticker Ticker) (*FREDSeriesInfo, error)
}

// HTTPClient defines the interface for HTTP operations.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// QueryOptions provides optional parameters for FRED API queries.
type QueryOptions struct {
	StartDate string
	EndDate   string
	Limit     int
	SortOrder string
}

// client implements the Client interface.
type client struct {
	apiKey     string
	httpClient HTTPClient
	baseURL    string
}

// NewClient creates a new FRED API client.
func NewClient(apiKey string) Client {
	return &client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: BaseURL,
	}
}

// NewClientWithHTTP creates a client with a custom HTTP client (for testing).
func NewClientWithHTTP(apiKey string, httpClient HTTPClient) Client {
	return &client{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    BaseURL,
	}
}

// GetSeriesObservations retrieves historical data for a ticker.
func (c *client) GetSeriesObservations(ctx context.Context, ticker Ticker, opts *QueryOptions) (*SeriesData, error) {
	if opts == nil {
		opts = &QueryOptions{
			Limit:     DefaultLimit,
			SortOrder: "desc",
		}
	}

	// Fetch observations
	apiURL := c.buildObservationsURL(ticker, opts)
	resp, err := c.doRequest(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch observations for %s: %w", ticker, err)
	}

	fredResp, err := c.parseObservationsResponse(resp)
	if err != nil {
		return nil, err
	}

	// Fetch series metadata for title, notes, etc.
	seriesInfo, err := c.GetSeriesInfo(ctx, ticker)
	if err != nil {
		// Don't fail if metadata fetch fails, just log it
		seriesInfo = &FREDSeriesInfo{
			Title:      ticker.Description(),
			Units:      fredResp.Units,
			UnitsShort: fredResp.UnitsShort,
			Frequency:  fredResp.Frequency,
		}
	}

	return &SeriesData{
		Ticker:       ticker,
		Description:  ticker.Description(),
		Title:        seriesInfo.Title,
		Observations: fredResp.Observations,
		Units:        seriesInfo.Units,
		UnitsShort:   seriesInfo.UnitsShort,
		Frequency:    seriesInfo.Frequency,
		Notes:        seriesInfo.Notes,
		LastUpdated:  time.Now(),
	}, nil
}

// GetSeriesInfo retrieves metadata for a ticker.
func (c *client) GetSeriesInfo(ctx context.Context, ticker Ticker) (*FREDSeriesInfo, error) {
	apiURL := c.buildSeriesURL(ticker)

	resp, err := c.doRequest(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series info for %s: %w", ticker, err)
	}

	seriesResp, err := c.parseSeriesResponse(resp)
	if err != nil {
		return nil, err
	}

	if len(seriesResp.Seriess) == 0 {
		return nil, fmt.Errorf("no series info found for %s", ticker)
	}

	return &seriesResp.Seriess[0], nil
}

// GetLatestValue retrieves the most recent value for a ticker.
func (c *client) GetLatestValue(ctx context.Context, ticker Ticker) (*LatestValue, error) {
	opts := &QueryOptions{
		Limit:     1,
		SortOrder: "desc",
	}

	seriesData, err := c.GetSeriesObservations(ctx, ticker, opts)
	if err != nil {
		return nil, err
	}

	if len(seriesData.Observations) == 0 {
		return nil, fmt.Errorf("no observations found for %s", ticker)
	}

	latest := seriesData.Observations[0]
	return &LatestValue{
		Ticker:      ticker,
		Description: ticker.Description(),
		Value:       latest.Value,
		Date:        latest.Date,
		UpdatedAt:   time.Now(),
	}, nil
}

// GetMultipleLatest retrieves the latest values for multiple tickers.
func (c *client) GetMultipleLatest(ctx context.Context, tickers []Ticker) (*MultiTickerResponse, error) {
	results := make([]LatestValue, 0, len(tickers))

	for _, ticker := range tickers {
		latest, err := c.GetLatestValue(ctx, ticker)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest for %s: %w", ticker, err)
		}
		results = append(results, *latest)
	}

	return &MultiTickerResponse{
		Data:      results,
		Timestamp: time.Now(),
	}, nil
}

// buildObservationsURL constructs the API URL with query parameters.
func (c *client) buildObservationsURL(ticker Ticker, opts *QueryOptions) string {
	params := url.Values{}
	params.Add("series_id", ticker.String())
	params.Add("api_key", c.apiKey)
	params.Add("file_type", "json")

	if opts.StartDate != "" {
		params.Add("observation_start", opts.StartDate)
	}
	if opts.EndDate != "" {
		params.Add("observation_end", opts.EndDate)
	}
	if opts.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.SortOrder != "" {
		params.Add("sort_order", opts.SortOrder)
	}

	return fmt.Sprintf("%s/series/observations?%s", c.baseURL, params.Encode())
}

// buildSeriesURL constructs the URL for fetching series metadata.
func (c *client) buildSeriesURL(ticker Ticker) string {
	params := url.Values{}
	params.Add("series_id", ticker.String())
	params.Add("api_key", c.apiKey)
	params.Add("file_type", "json")

	return fmt.Sprintf("%s/series?%s", c.baseURL, params.Encode())
}

// doRequest performs an HTTP request with context.
func (c *client) doRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// parseObservationsResponse parses the FRED API observations JSON response.
func (c *client) parseObservationsResponse(resp *http.Response) (*FREDAPIResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var fredResp FREDAPIResponse
	if err := json.Unmarshal(body, &fredResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &fredResp, nil
}

// parseSeriesResponse parses the FRED API series metadata JSON response.
func (c *client) parseSeriesResponse(resp *http.Response) (*FREDSeriesResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var seriesResp FREDSeriesResponse
	if err := json.Unmarshal(body, &seriesResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &seriesResp, nil
}
