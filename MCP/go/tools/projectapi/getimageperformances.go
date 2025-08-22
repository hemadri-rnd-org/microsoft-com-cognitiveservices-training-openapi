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

func GetimageperformancesHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if val, ok := args["tagIds"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("tagIds=%v", val))
		}
		if val, ok := args["orderBy"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("orderBy=%v", val))
		}
		if val, ok := args["take"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("take=%v", val))
		}
		if val, ok := args["skip"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("skip=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/projects/%s/iterations/%s/performance/images%s", cfg.BaseURL, projectId, iterationId, queryString)
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
		var result []ImagePerformance
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

func CreateGetimageperformancesTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_projects_projectId_iterations_iterationId_performance_images",
		mcp.WithDescription("Get image with its prediction for a given project iteration."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The project id.")),
		mcp.WithString("iterationId", mcp.Required(), mcp.Description("The iteration id. Defaults to workspace.")),
		mcp.WithArray("tagIds", mcp.Description("A list of tags ids to filter the images. Defaults to all tagged images when null. Limited to 20.")),
		mcp.WithString("orderBy", mcp.Description("The ordering. Defaults to newest.")),
		mcp.WithNumber("take", mcp.Description("Maximum number of images to return. Defaults to 50, limited to 256.")),
		mcp.WithNumber("skip", mcp.Description("Number of images to skip before beginning the image batch. Defaults to 0.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetimageperformancesHandler(cfg),
	}
}
