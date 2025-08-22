package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/custom-vision-training-client/mcp-server/config"
	"github.com/custom-vision-training-client/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetimageregionproposalsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		imageIdVal, ok := args["imageId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: imageId"), nil
		}
		imageId, ok := imageIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: imageId"), nil
		}
		url := fmt.Sprintf("%s/projects/%s/images/%s/regionproposals", cfg.BaseURL, projectId, imageId)
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
		var result models.ImageRegionProposal
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

func CreateGetimageregionproposalsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_projects_projectId_images_imageId_regionproposals",
		mcp.WithDescription("Get region proposals for an image. Returns empty array if no proposals are found."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The project id.")),
		mcp.WithString("imageId", mcp.Required(), mcp.Description("The image id.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetimageregionproposalsHandler(cfg),
	}
}
