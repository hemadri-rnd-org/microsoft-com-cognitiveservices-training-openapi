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

func PublishiterationHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if val, ok := args["publishName"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("publishName=%v", val))
		}
		if val, ok := args["predictionId"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("predictionId=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/projects/%s/iterations/%s/publish%s", cfg.BaseURL, projectId, iterationId, queryString)
		req, err := http.NewRequest("POST", url, nil)
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
		var result bool
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

func CreatePublishiterationTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_projects_projectId_iterations_iterationId_publish",
		mcp.WithDescription("Publish a specific iteration."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The project id.")),
		mcp.WithString("iterationId", mcp.Required(), mcp.Description("The iteration id.")),
		mcp.WithString("publishName", mcp.Required(), mcp.Description("The name to give the published iteration.")),
		mcp.WithString("predictionId", mcp.Required(), mcp.Description("The id of the prediction resource to publish to.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    PublishiterationHandler(cfg),
	}
}
