package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/custom-vision-training-client/mcp-server/config"
	"github.com/custom-vision-training-client/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetiterationperformanceHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		projectIdVal, ok := args["projectId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: projectId"), nil
		}
		projectId, ok := projectIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: projectId"), nil
		}
		iterationIdVal, ok := args["iterationId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: iterationId"), nil
		}
		iterationId, ok := iterationIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: iterationId"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["threshold"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("threshold=%v", val))
		}
		if val, ok := args["overlapThreshold"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("overlapThreshold=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/projects/%s/iterations/%s/performance%s", cfg.BaseURL, projectId, iterationId, queryString)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			req.Header.Set("Training-Key", cfg.APIKey)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.IterationPerformance
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetiterationperformanceTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_projects_projectId_iterations_iterationId_performance",
		mcp.WithDescription("Get detailed performance information about an iteration."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The id of the project the iteration belongs to.")),
		mcp.WithString("iterationId", mcp.Required(), mcp.Description("The id of the iteration to get.")),
		mcp.WithString("threshold", mcp.Description("The threshold used to determine true predictions.")),
		mcp.WithString("overlapThreshold", mcp.Description("If applicable, the bounding box overlap threshold used to determine true predictions.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetiterationperformanceHandler(cfg),
	}
}
