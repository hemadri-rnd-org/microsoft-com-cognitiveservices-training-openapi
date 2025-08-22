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

func SuggesttagsandregionsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		queryParams := make([]string, 0)
		if val, ok := args["iterationId"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("iterationId=%v", val))
		}
		if val, ok := args["imageIds"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("imageIds=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/projects/%s/tagsandregions/suggestions%s", cfg.BaseURL, projectId, queryString)
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
		var result []SuggestedTagAndRegion
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

func CreateSuggesttagsandregionsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_projects_projectId_tagsandregions_suggestions",
		mcp.WithDescription("Suggest tags and regions for an array/batch of untagged images. Returns empty array if no tags are found."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The project id.")),
		mcp.WithString("iterationId", mcp.Required(), mcp.Description("IterationId to use for tag and region suggestion.")),
		mcp.WithArray("imageIds", mcp.Required(), mcp.Description("Array of image ids tag suggestion are needed for. Use GetUntaggedImages API to get imageIds.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    SuggesttagsandregionsHandler(cfg),
	}
}
