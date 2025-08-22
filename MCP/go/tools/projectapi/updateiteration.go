package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/custom-vision-training-client/mcp-server/config"
	"github.com/custom-vision-training-client/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func UpdateiterationHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		// Create properly typed request body using the generated schema
		var requestBody models.Iteration
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/projects/%s/iterations/%s", cfg.BaseURL, projectId, iterationId)
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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
		var result models.Iteration
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

func CreateUpdateiterationTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("patch_projects_projectId_iterations_iterationId",
		mcp.WithDescription("Update a specific iteration."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("Project id.")),
		mcp.WithString("iterationId", mcp.Required(), mcp.Description("Iteration id.")),
		mcp.WithNumber("reservedBudgetInHours", mcp.Description("Input parameter: Gets the reserved advanced training budget for the iteration.")),
		mcp.WithString("status", mcp.Description("Input parameter: Gets the current iteration status.")),
		mcp.WithArray("exportableTo", mcp.Description("Input parameter: A set of platforms this iteration can export to.")),
		mcp.WithNumber("trainingTimeInMinutes", mcp.Description("Input parameter: Gets the training time for the iteration.")),
		mcp.WithString("classificationType", mcp.Description("Input parameter: Gets the classification type of the project.")),
		mcp.WithString("originalPublishResourceId", mcp.Description("Input parameter: Resource Provider Id this iteration was originally published to.")),
		mcp.WithString("created", mcp.Description("Input parameter: Gets the time this iteration was completed.")),
		mcp.WithString("publishName", mcp.Description("Input parameter: Name of the published model.")),
		mcp.WithString("id", mcp.Description("Input parameter: Gets the id of the iteration.")),
		mcp.WithString("domainId", mcp.Description("Input parameter: Get or sets a guid of the domain the iteration has been trained on.")),
		mcp.WithBoolean("exportable", mcp.Description("Input parameter: Whether the iteration can be exported to another format for download.")),
		mcp.WithString("trainedAt", mcp.Description("Input parameter: Gets the time this iteration was last modified.")),
		mcp.WithString("projectId", mcp.Description("Input parameter: Gets the project id of the iteration.")),
		mcp.WithString("lastModified", mcp.Description("Input parameter: Gets the time this iteration was last modified.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Input parameter: Gets or sets the name of the iteration.")),
		mcp.WithString("trainingType", mcp.Description("Input parameter: Gets the training type of the iteration.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    UpdateiterationHandler(cfg),
	}
}
