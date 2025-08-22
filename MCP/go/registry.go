package main

import (
	"github.com/custom-vision-training-client/mcp-server/config"
	"github.com/custom-vision-training-client/mcp-server/models"
	tools_imageapi "github.com/custom-vision-training-client/mcp-server/tools/imageapi"
	tools_domainsapi "github.com/custom-vision-training-client/mcp-server/tools/domainsapi"
	tools_projectapi "github.com/custom-vision-training-client/mcp-server/tools/projectapi"
	tools_tagsapi "github.com/custom-vision-training-client/mcp-server/tools/tagsapi"
	tools_suggesttagsandregionsapi "github.com/custom-vision-training-client/mcp-server/tools/suggesttagsandregionsapi"
	tools_predictionsapi "github.com/custom-vision-training-client/mcp-server/tools/predictionsapi"
	tools_imageregionproposalapi "github.com/custom-vision-training-client/mcp-server/tools/imageregionproposalapi"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_imageapi.CreateGetimagesbyidsTool(cfg),
		tools_domainsapi.CreateGetdomainTool(cfg),
		tools_projectapi.CreateGetimageperformancesTool(cfg),
		tools_imageapi.CreateDeleteimageregionsTool(cfg),
		tools_imageapi.CreateCreateimageregionsTool(cfg),
		tools_tagsapi.CreateGettagsTool(cfg),
		tools_tagsapi.CreateCreatetagTool(cfg),
		tools_imageapi.CreateGetuntaggedimagecountTool(cfg),
		tools_projectapi.CreateGetprojectsTool(cfg),
		tools_projectapi.CreateCreateprojectTool(cfg),
		tools_imageapi.CreateCreateimagesfromfilesTool(cfg),
		tools_imageapi.CreateDeleteimagetagsTool(cfg),
		tools_imageapi.CreateCreateimagetagsTool(cfg),
		tools_projectapi.CreateExportprojectTool(cfg),
		tools_imageapi.CreateGettaggedimagecountTool(cfg),
		tools_suggesttagsandregionsapi.CreateSuggesttagsandregionsTool(cfg),
		tools_imageapi.CreateCreateimagesfrompredictionsTool(cfg),
		tools_projectapi.CreateGetiterationsTool(cfg),
		tools_imageapi.CreateCreateimagesfromurlsTool(cfg),
		tools_imageapi.CreateGettaggedimagesTool(cfg),
		tools_projectapi.CreateTrainprojectTool(cfg),
		tools_imageapi.CreateGetuntaggedimagesTool(cfg),
		tools_predictionsapi.CreateDeletepredictionTool(cfg),
		tools_imageapi.CreateQuerysuggestedimagesTool(cfg),
		tools_imageregionproposalapi.CreateGetimageregionproposalsTool(cfg),
		tools_projectapi.CreateExportiterationTool(cfg),
		tools_projectapi.CreateGetexportsTool(cfg),
		tools_projectapi.CreateImportprojectTool(cfg),
		tools_imageapi.CreateDeleteimagesTool(cfg),
		tools_projectapi.CreateGetiterationperformanceTool(cfg),
		tools_tagsapi.CreateUpdatetagTool(cfg),
		tools_tagsapi.CreateDeletetagTool(cfg),
		tools_tagsapi.CreateGettagTool(cfg),
		tools_projectapi.CreateDeleteiterationTool(cfg),
		tools_projectapi.CreateGetiterationTool(cfg),
		tools_projectapi.CreateUpdateiterationTool(cfg),
		tools_projectapi.CreateGetimageperformancecountTool(cfg),
		tools_imageapi.CreateQuerysuggestedimagecountTool(cfg),
		tools_projectapi.CreateUnpublishiterationTool(cfg),
		tools_projectapi.CreatePublishiterationTool(cfg),
		tools_predictionsapi.CreateQuerypredictionsTool(cfg),
		tools_domainsapi.CreateGetdomainsTool(cfg),
		tools_projectapi.CreateUpdateprojectTool(cfg),
		tools_projectapi.CreateDeleteprojectTool(cfg),
		tools_projectapi.CreateGetprojectTool(cfg),
		tools_predictionsapi.CreateQuicktestimageurlTool(cfg),
	}
}
