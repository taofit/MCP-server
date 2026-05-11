package main

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	NWSAPIBase = "https://api.weather.gov"
	UserAgent  = "weather-app/1.0"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Weather",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		Instructions: "You are a helpful assistant specialized in US weather. Your goal is to provide accurate weather forecast information for any location within the United States. If a user asks about a location outside the US, you must politely decline, explaining that your capabilities are limited to American soil. Always respond in the language the user is speaking.",
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get weather forecast for a location in United States. If a user asks for a location outside the US, explain that this tool only covers American soil.",
	}, getForecast)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_alerts",
		Description: "Get weather alerts for a US state. If a user asks for a location outside the US, explain that this tool only covers American soil.",
	}, getAlerts)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func makeNWSRequest[T any](ctx context.Context, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Add("Accept", "application/geo+json")
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/geo+json")

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to make request to %s, %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return &result, nil
}

func formatAlert(alert AlertFeature) string {
	props := alert.Properties
	event := cmp.Or(props.Event, "Unknown")
	areaDesc := cmp.Or(props.AreaDesc, "Unknown")
	severity := cmp.Or(props.Severity, "Unknown")
	description := cmp.Or(props.Description, "No description available")
	instruction := cmp.Or(props.Instruction, "No specific instructions provided")

	return fmt.Sprintf(`
	Event: %s
	Area: %s
	Severity: %s
	Description: %s
	Instructions: %s
	`, event, areaDesc, severity, description, instruction)
}

func formatPeriod(period ForecastPeriod) string {
	return fmt.Sprintf(`
	%s:
	Temperature: %d°%s
	Wind: %s %s
	Forecast: %s
	`, period.Name, period.Temperature, period.TemperatureUnit, period.WindSpeed, period.WindDirection, period.DetailedForecast)
}

var getForecast mcp.ToolHandlerFor[ForecastInput, any] = func(ctx context.Context, req *mcp.CallToolRequest, input ForecastInput) (*mcp.CallToolResult, any, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	pointsURL := fmt.Sprintf("%s/points/%f,%f", NWSAPIBase, input.Latitude, input.Longitude)
	pointsData, err := makeNWSRequest[PointsResponse](ctx, pointsURL)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unable to fetch forecast data for this location: %v", err)},
			},
		}, nil, nil
	}

	forecastURL := pointsData.Properties.Forecast
	if forecastURL == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Unable to fetch forecast URL."},
			},
		}, nil, nil
	}

	forecastData, err := makeNWSRequest[ForecastResponse](ctx, forecastURL)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unable to fetch detailedforecast: %v", err)},
			},
		}, nil, nil
	}

	periods := forecastData.Properties.Periods
	if len(periods) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No forecast periods available."},
			},
		}, nil, nil
	}

	var forecasts []string
	for i := range min(5, len(periods)) {
		forecasts = append(forecasts, formatPeriod(periods[i]))
	}

	result := strings.Join(forecasts, "\n---\n")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

var getAlerts mcp.ToolHandlerFor[AlertsInput, any] = func(ctx context.Context, req *mcp.CallToolRequest, input AlertsInput) (*mcp.CallToolResult, any, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	stateCode := strings.ToUpper(input.State)
	alertsURL := fmt.Sprintf("%s/alerts/active/area/%s", NWSAPIBase, stateCode)

	alertsData, err := makeNWSRequest[AlertsResponse](ctx, alertsURL)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unable to fetch alerts data or no alerts found: %v", err)},
			},
		}, nil, nil
	}

	if len(alertsData.Features) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No active alerts for this state."},
			},
		}, nil, nil
	}

	var alerts []string
	for _, feature := range alertsData.Features {
		alerts = append(alerts, formatAlert(feature))
	}

	result := strings.Join(alerts, "\n---\n")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}
