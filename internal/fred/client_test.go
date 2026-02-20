package fred

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockHTTPClient implements HTTPClient for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

// TestNewClient verifies client initialization.
func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

// TestNewClientWithHTTP verifies client initialization with custom HTTP client.
func TestNewClientWithHTTP(t *testing.T) {
	apiKey := "test-api-key"
	mockHTTP := &MockHTTPClient{}

	client := NewClientWithHTTP(apiKey, mockHTTP)

	if client == nil {
		t.Fatal("NewClientWithHTTP returned nil")
	}
}

// TestGetLatestValue verifies fetching the latest value for a ticker.
func TestGetLatestValue(t *testing.T) {
	mockResp := FREDAPIResponse{
		Observations: []Observation{
			{Date: "2024-01-15", Value: "50000.5"},
		},
	}

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(mockResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	result, err := client.GetLatestValue(ctx, TickerWALCL)
	if err != nil {
		t.Fatalf("GetLatestValue failed: %v", err)
	}

	if result.Ticker != TickerWALCL {
		t.Errorf("Expected ticker %s, got %s", TickerWALCL, result.Ticker)
	}

	if result.Value != "50000.5" {
		t.Errorf("Expected value 50000.5, got %s", result.Value)
	}

	if result.Date != "2024-01-15" {
		t.Errorf("Expected date 2024-01-15, got %s", result.Date)
	}

	if result.Description == "" {
		t.Error("Description should not be empty")
	}
}

// TestGetLatestValueNoObservations verifies error handling for empty response.
func TestGetLatestValueNoObservations(t *testing.T) {
	mockResp := FREDAPIResponse{
		Observations: []Observation{},
	}

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(mockResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	_, err := client.GetLatestValue(ctx, TickerWALCL)
	if err == nil {
		t.Error("Expected error for empty observations, got nil")
	}
}

// TestGetSeriesObservations verifies fetching historical data.
func TestGetSeriesObservations(t *testing.T) {
	callCount := 0
	
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			
			// First call: observations
			if callCount == 1 {
				mockResp := FREDAPIResponse{
					Observations: []Observation{
						{Date: "2024-01-15", Value: "100.5"},
						{Date: "2024-01-14", Value: "99.8"},
					},
					Units:      "Millions of Dollars",
					UnitsShort: "Mil. of $",
					Frequency:  "Weekly",
				}
				body, _ := json.Marshal(mockResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			
			// Second call: series info
			mockSeriesResp := FREDSeriesResponse{
				Seriess: []FREDSeriesInfo{
					{
						ID:             "CPIAUCSL",
						Title:          "Consumer Price Index for All Urban Consumers",
						Units:          "Index 1982-1984=100",
						UnitsShort:     "Index",
						Frequency:      "Monthly",
						FrequencyShort: "M",
						Notes:          "Consumer Price Index data",
					},
				},
			}
			body, _ := json.Marshal(mockSeriesResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	opts := &QueryOptions{
		Limit:     10,
		SortOrder: "desc",
	}

	result, err := client.GetSeriesObservations(ctx, TickerCPIAUCSL, opts)
	if err != nil {
		t.Fatalf("GetSeriesObservations failed: %v", err)
	}

	if len(result.Observations) != 2 {
		t.Errorf("Expected 2 observations, got %d", len(result.Observations))
	}

	if result.Ticker != TickerCPIAUCSL {
		t.Errorf("Expected ticker %s, got %s", TickerCPIAUCSL, result.Ticker)
	}
	
	if result.Title == "" {
		t.Error("Title should not be empty")
	}
	
	if result.Units == "" {
		t.Error("Units should not be empty")
	}
	
	if result.UnitsShort == "" {
		t.Error("UnitsShort should not be empty")
	}
	
	if result.Frequency == "" {
		t.Error("Frequency should not be empty")
	}
}

// TestGetSeriesObservationsWithNilOptions verifies default options.
func TestGetSeriesObservationsWithNilOptions(t *testing.T) {
	callCount := 0
	
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			
			if callCount == 1 {
				mockResp := FREDAPIResponse{
					Observations: []Observation{
						{Date: "2024-01-15", Value: "100.5"},
					},
				}
				body, _ := json.Marshal(mockResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			
			// Series info call
			mockSeriesResp := FREDSeriesResponse{
				Seriess: []FREDSeriesInfo{
					{
						Title:      "Test Series",
						Units:      "Test Units",
						UnitsShort: "TU",
						Frequency:  "Daily",
					},
				},
			}
			body, _ := json.Marshal(mockSeriesResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	result, err := client.GetSeriesObservations(ctx, TickerWALCL, nil)
	if err != nil {
		t.Fatalf("GetSeriesObservations with nil options failed: %v", err)
	}

	if len(result.Observations) == 0 {
		t.Error("Expected observations, got empty")
	}
}

// TestGetMultipleLatest verifies fetching multiple tickers.
func TestGetMultipleLatest(t *testing.T) {
	callCount := 0
	tickers := []string{"50000.5", "3000.2", "0.05"}

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Each ticker makes 2 calls: observations + series info
			tickerIndex := callCount / 2
			isSeriesInfo := callCount%2 == 1
			
			if isSeriesInfo {
				// Series info call
				mockSeriesResp := FREDSeriesResponse{
					Seriess: []FREDSeriesInfo{
						{
							Title:      "Test Series",
							Units:      "Test Units",
							UnitsShort: "TU",
							Frequency:  "Daily",
						},
					},
				}
				body, _ := json.Marshal(mockSeriesResp)
				callCount++
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			
			// Observations call
			mockResp := FREDAPIResponse{
				Observations: []Observation{
					{Date: "2024-01-15", Value: tickers[tickerIndex]},
				},
			}
			body, _ := json.Marshal(mockResp)
			callCount++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	tickerList := []Ticker{TickerWALCL, TickerCPIAUCSL, TickerFEDFUNDS}
	result, err := client.GetMultipleLatest(ctx, tickerList)
	if err != nil {
		t.Fatalf("GetMultipleLatest failed: %v", err)
	}

	if len(result.Data) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Data))
	}

	for i, data := range result.Data {
		if data.Ticker != tickerList[i] {
			t.Errorf("Expected ticker %s at index %d, got %s", tickerList[i], i, data.Ticker)
		}
	}
}

// TestGetMultipleLatestWithError verifies error handling.
func TestGetMultipleLatestWithError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	tickerList := []Ticker{TickerWALCL}
	_, err := client.GetMultipleLatest(ctx, tickerList)
	if err == nil {
		t.Error("Expected error for network failure, got nil")
	}
}

// TestBuildObservationsURL verifies URL construction.
func TestBuildObservationsURL(t *testing.T) {
	c := &client{
		apiKey:  "test-key",
		baseURL: BaseURL,
	}

	opts := &QueryOptions{
		StartDate: "2024-01-01",
		EndDate:   "2024-01-31",
		Limit:     50,
		SortOrder: "asc",
	}

	url := c.buildObservationsURL(TickerWALCL, opts)

	if url == "" {
		t.Fatal("buildObservationsURL returned empty string")
	}

	if !contains(url, "series_id=WALCL") {
		t.Error("URL should contain series_id parameter")
	}

	if !contains(url, "api_key=test-key") {
		t.Error("URL should contain api_key parameter")
	}

	if !contains(url, "observation_start=2024-01-01") {
		t.Error("URL should contain start date parameter")
	}

	if !contains(url, "observation_end=2024-01-31") {
		t.Error("URL should contain end date parameter")
	}

	if !contains(url, "limit=50") {
		t.Error("URL should contain limit parameter")
	}

	if !contains(url, "sort_order=asc") {
		t.Error("URL should contain sort_order parameter")
	}
}

// TestDoRequestWithContext verifies context handling.
func TestDoRequestWithContext(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Context() == nil {
				t.Error("Request context should not be nil")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			}, nil
		},
	}

	c := &client{
		apiKey:     "test-key",
		httpClient: mockHTTP,
		baseURL:    BaseURL,
	}

	ctx := context.Background()
	url := fmt.Sprintf("%s/test", BaseURL)

	_, err := c.doRequest(ctx, url)
	if err != nil {
		t.Errorf("doRequest failed: %v", err)
	}
}

// TestDoRequestWithHTTPError verifies HTTP error handling.
func TestDoRequestWithHTTPError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader([]byte("Bad Request"))),
			}, nil
		},
	}

	c := &client{
		apiKey:     "test-key",
		httpClient: mockHTTP,
		baseURL:    BaseURL,
	}

	ctx := context.Background()
	url := fmt.Sprintf("%s/test", BaseURL)

	_, err := c.doRequest(ctx, url)
	if err == nil {
		t.Error("Expected error for HTTP 400, got nil")
	}
}

// TestDoRequestWithNetworkError verifies network error handling.
func TestDoRequestWithNetworkError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("connection timeout")
		},
	}

	c := &client{
		apiKey:     "test-key",
		httpClient: mockHTTP,
		baseURL:    BaseURL,
	}

	ctx := context.Background()
	url := fmt.Sprintf("%s/test", BaseURL)

	_, err := c.doRequest(ctx, url)
	if err == nil {
		t.Error("Expected error for network failure, got nil")
	}
}

// TestParseResponse verifies JSON parsing.
func TestParseResponse(t *testing.T) {
	mockResp := FREDAPIResponse{
		RealtimeStart: "2024-01-01",
		RealtimeEnd:   "2024-01-31",
		Observations: []Observation{
			{Date: "2024-01-15", Value: "100.5"},
		},
		Count:      1,
		Units:      "Millions of Dollars",
		UnitsShort: "Mil. of $",
		Frequency:  "Weekly",
	}

	body, _ := json.Marshal(mockResp)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	c := &client{}
	result, err := c.parseObservationsResponse(resp)

	if err != nil {
		t.Fatalf("parseObservationsResponse failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Expected count 1, got %d", result.Count)
	}

	if len(result.Observations) != 1 {
		t.Errorf("Expected 1 observation, got %d", len(result.Observations))
	}
	
	if result.Units != "Millions of Dollars" {
		t.Errorf("Expected units 'Millions of Dollars', got %s", result.Units)
	}
	
	if result.UnitsShort != "Mil. of $" {
		t.Errorf("Expected units_short 'Mil. of $', got %s", result.UnitsShort)
	}
}

// TestParseResponseWithInvalidJSON verifies JSON error handling.
func TestParseResponseWithInvalidJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
	}

	c := &client{}
	_, err := c.parseObservationsResponse(resp)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestGetSeriesInfo verifies series metadata retrieval.
func TestGetSeriesInfo(t *testing.T) {
	mockSeriesResp := FREDSeriesResponse{
		Seriess: []FREDSeriesInfo{
			{
				ID:             "WALCL",
				Title:          "Federal Reserve Total Assets",
				Units:          "Millions of Dollars",
				UnitsShort:     "Mil. of $",
				Frequency:      "Weekly, As of Wednesday",
				FrequencyShort: "W",
				Notes:          "Assets: Total Assets: Total Assets (Less Eliminations from Consolidation): Wednesday Level",
			},
		},
	}

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(mockSeriesResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	result, err := client.GetSeriesInfo(ctx, TickerWALCL)
	if err != nil {
		t.Fatalf("GetSeriesInfo failed: %v", err)
	}

	if result.ID != "WALCL" {
		t.Errorf("Expected ID WALCL, got %s", result.ID)
	}

	if result.Title != "Federal Reserve Total Assets" {
		t.Errorf("Expected title 'Federal Reserve Total Assets', got %s", result.Title)
	}

	if result.Units != "Millions of Dollars" {
		t.Errorf("Expected units 'Millions of Dollars', got %s", result.Units)
	}

	if result.UnitsShort != "Mil. of $" {
		t.Errorf("Expected units_short 'Mil. of $', got %s", result.UnitsShort)
	}

	if result.Frequency == "" {
		t.Error("Frequency should not be empty")
	}

	if result.Notes == "" {
		t.Error("Notes should not be empty")
	}
}

// TestGetSeriesInfoNoResults verifies error handling when no series found.
func TestGetSeriesInfoNoResults(t *testing.T) {
	mockSeriesResp := FREDSeriesResponse{
		Seriess: []FREDSeriesInfo{},
	}

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(mockSeriesResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx := context.Background()

	_, err := client.GetSeriesInfo(ctx, TickerWALCL)
	if err == nil {
		t.Error("Expected error for empty series response, got nil")
	}
}

// TestBuildSeriesURL verifies series URL construction.
func TestBuildSeriesURL(t *testing.T) {
	c := &client{
		apiKey:  "test-key",
		baseURL: BaseURL,
	}

	url := c.buildSeriesURL(TickerWALCL)

	if url == "" {
		t.Fatal("buildSeriesURL returned empty string")
	}

	if !contains(url, "series_id=WALCL") {
		t.Error("URL should contain series_id parameter")
	}

	if !contains(url, "api_key=test-key") {
		t.Error("URL should contain api_key parameter")
	}

	if !contains(url, "/series?") {
		t.Error("URL should contain /series endpoint")
	}
}

// TestContextCancellation verifies context cancellation is respected.
func TestContextCancellation(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			time.Sleep(100 * time.Millisecond)
			return nil, req.Context().Err()
		},
	}

	client := NewClientWithHTTP("test-key", mockHTTP)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetLatestValue(ctx, TickerWALCL)
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
}

// Helper function to check if string contains substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
