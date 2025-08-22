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

func CreateprojectHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["name"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("name=%v", val))
		}
		if val, ok := args["description"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("description=%v", val))
		}
		if val, ok := args["domainId"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("domainId=%v", val))
		}
		if val, ok := args["classificationType"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("classificationType=%v", val))
		}
		if val, ok := args["targetExportPlatforms"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("targetExportPlatforms=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/projects%s", cfg.BaseURL, queryString)
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
		var result models.Project
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

func CreateCreateprojectTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_projects",
		mcp.WithDescription("Create a project."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the project.")),
		mcp.WithString("description", mcp.Description("The description of the project.")),
		mcp.WithString("domainId", mcp.Description("The id of the domain to use for this project. Defaults to General.")),
		mcp.WithString("classificationType", mcp.Description("The type of classifier to create for this project.")),
		mcp.WithArray("targetExportPlatforms", mcp.Description("List of platforms the trained model is intending exporting to.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    CreateprojectHandler(cfg),
	}
}
