package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"bytes"

	"github.com/custom-vision-training-client/mcp-server/config"
	"github.com/custom-vision-training-client/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func TrainprojectHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if val, ok := args["trainingType"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("trainingType=%v", val))
		}
		if val, ok := args["reservedBudgetInHours"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("reservedBudgetInHours=%v", val))
		}
		if val, ok := args["forceTrain"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("forceTrain=%v", val))
		}
		if val, ok := args["notificationEmailAddress"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("notificationEmailAddress=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		// Create properly typed request body using the generated schema
		var requestBody models.TrainingParameters
		
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
		url := fmt.Sprintf("%s/projects/%s/train%s", cfg.BaseURL, projectId, queryString)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
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

func CreateTrainprojectTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_projects_projectId_train",
		mcp.WithDescription("Queues project for training."),
		mcp.WithString("projectId", mcp.Required(), mcp.Description("The project id.")),
		mcp.WithString("trainingType", mcp.Description("The type of training to use to train the project (default: Regular).")),
		mcp.WithNumber("reservedBudgetInHours", mcp.Description("The number of hours reserved as budget for training (if applicable).")),
		mcp.WithBoolean("forceTrain", mcp.Description("Whether to force train even if dataset and configuration does not change (default: false).")),
		mcp.WithString("notificationEmailAddress", mcp.Description("The email address to send notification to when training finishes (default: null).")),
		mcp.WithArray("selectedTags", mcp.Description("Input parameter: List of tags selected for this training session, other tags in the project will be ignored.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    TrainprojectHandler(cfg),
	}
}
