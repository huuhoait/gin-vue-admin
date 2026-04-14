package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GVAReviewer GVAcode review tool
type GVAReviewer struct{}

// init register tool
func init() {
	RegisterTool(&GVAReviewer{})
}

// ReviewRequest reviewCheckRequestStructure
type ReviewRequest struct {
	UserRequirement string   `json:"userRequirement"` // ViaPassrequirement_analyzeAfterofuser requirement
	GeneratedFiles  []string `json:"generatedFiles"`  // gva_executeCreateofFileList
}

// ReviewResponse reviewCheckResponseStructure
type ReviewResponse struct {
	Success          bool   `json:"success"`          // YesNoreviewChecksucceeded
	Message          string `json:"message"`          // reviewCheckResultMessage
	AdjustmentPrompt string `json:"adjustmentPrompt"` // adjustArrangeGenerationCodeofPrompt
	ReviewDetails    string `json:"reviewDetails"`    // DetailsofreviewCheckResult
}

// New CreateGVAcode review tool
func (g *GVAReviewer) New() mcp.Tool {
	return mcp.NewTool("gva_review",
		mcp.WithDescription(`**GVAcode review tool - Atgva_executeInvokeAfterUse**

**Core features:**
- ReceiveViaPassrequirement_analyzeHandleofuser requirementAndgva_executeGenerateofFileList
- PartAnalyzeGenerateofGenerationCodeYesNoSatisfyUserofOriginalNeedRequest
- CheckYesNoInvolveToAssociation, Exchangemutualetc.ComplexFunction
- IfGenerationCodeNotSatisfyNeedRequest, ProvideadjustArrangeSuggestAndNewofprompt

**Use cases:**
- Atgva_executesucceededExecuteAfterInvoke
- Used forVerifyGenerateofGenerationCodeYesNoCompleteArrangeSatisfyuser requirement
- CheckModuleIntervalofassociationYesNoCorrectImplement
- DiscoverMissingofExchangemutualFunctionOrBusinessLogic

**WorkflowFlow:**
1. ReceiveUserOriginalNeedRequestAndGenerateofFileList
2. PartAnalyzeNeedRequestInCloseKeyFunctionPoint
3. CheckGenerateofFileYesNoCoverAllFunction
4. IdentifyMissingofassociation, ExchangemutualFunctionetc.
5. GenerateadjustArrangeSuggestAndNewofDevelopmentprompt

**InputOutcontent:**
- reviewCheckResultAndYesNoNeedadjustArrange
- DetailsofMissingFunctionPartAnalyze
- ForNatureofGenerationCodeadjustArrangeSuggest
- CanDirectUseofDevelopmentprompt

**Important notes:**
- ThisUtilitySpecialDoorUsed forGenerationCodeQualityAmountreviewCheck, NotExecuteActualofGenerationCodeUpdate
- RepeatPointCloseNoteModuleIntervalAssociation, UserExchangemutual, BusinessWorkflowCompleteArrangeNature
- ProvideofadjustArrangeSuggestShouldThisToolBodyCanExecute`),
		mcp.WithString("userRequirement",
			mcp.Description("ViaPassrequirement_analyzeHandleAfterofuser requirementdescription, PackageIncludeDetailsofFunctionRequirementAndFieldInformation"),
			mcp.Required(),
		),
		mcp.WithString("generatedFiles",
			mcp.Description("gva_executeCreateofFileList, JSONStringFormat, PackageIncludeAllGenerateofAfterEndAndFrontendFilepath"),
			mcp.Required(),
		),
	)
}

// Handle HandlereviewCheckRequest
func (g *GVAReviewer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// getuser requirement
	userRequirementData, ok := request.GetArguments()["userRequirement"]
	if !ok {
		return nil, errors.New("Invalid parameters:userRequirement must be provided")
	}

	userRequirement, ok := userRequirementData.(string)
	if !ok {
		return nil, errors.New("Invalid parameters:userRequirement MustYesStringtype")
	}

	// getGenerateofFileList
	generatedFilesData, ok := request.GetArguments()["generatedFiles"]
	if !ok {
		return nil, errors.New("Invalid parameters:generatedFiles must be provided")
	}

	generatedFilesStr, ok := generatedFilesData.(string)
	if !ok {
		return nil, errors.New("Invalid parameters:generatedFiles MustYesJSONString")
	}

	// ParseJSONStringForStringArray
	var generatedFiles []string
	err := json.Unmarshal([]byte(generatedFilesStr), &generatedFiles)
	if err != nil {
		return nil, fmt.Errorf("ParsegeneratedFilesfailed: %v", err)
	}

	if len(generatedFiles) == 0 {
		return nil, errors.New("Invalid parameters:generatedFiles cannot be empty")
	}

	// DirectGenerateadjustArrangePrompt, NotPerformComplexPartAnalyze
	adjustmentPrompt := g.generateAdjustmentPrompt(userRequirement, generatedFiles)

	// BuildSimplifyofreviewCheckDetails
	reviewDetails := fmt.Sprintf("📋 **GenerationCodereviewCheckReport**\n\n **UserOriginalNeedRequest:**\n%s\n\n **AlreadyGenerateFileQuantity:** %d\n\n **SuggestPerformGenerationCodeOptimizationAndCompletegood**", userRequirement, len(generatedFiles))

	// BuildreviewCheckResult
	reviewResult := &ReviewResponse{
		Success:          true,
		Message:          "GenerationCodereviewCheckComplete",
		AdjustmentPrompt: adjustmentPrompt,
		ReviewDetails:    reviewDetails,
	}

	// serialize response
	responseJSON, err := json.MarshalIndent(reviewResult, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serialize review resultfailed: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("GenerationCodereviewCheckResult:\n\n%s", string(responseJSON))),
		},
	}, nil
}

// generateAdjustmentPrompt GenerateadjustArrangeGenerationCodeofPrompt
func (g *GVAReviewer) generateAdjustmentPrompt(userRequirement string, generatedFiles []string) string {
	var prompt strings.Builder

	prompt.WriteString("🔧 **GenerationCodeadjustArrangeInstructionGuide Prompt:**\n\n")
	prompt.WriteString(fmt.Sprintf("**UserofOriginalNeedRequestFor:** %s\n\n", userRequirement))
	prompt.WriteString("**ViaPassGVAGenerateAfterofFileHaveSuch AsDowncontent:**\n")
	for _, file := range generatedFiles {
		prompt.WriteString(fmt.Sprintf("- %s\n", file))
	}
	prompt.WriteString("\n")

	prompt.WriteString("**PleaseHelpIOptimizationAndCompletegoodGenerationCode, ensure:**\n")
	prompt.WriteString("1. GenerationCodeCompleteAllSatisfyUserofOriginalNeedRequest\n")
	prompt.WriteString("2. CompletegoodModuleIntervalofassociation, ensureDataConsistentNature\n")
	prompt.WriteString("3. ImplementAllmustRequireofUserExchangemutualFunction\n")
	prompt.WriteString("4. protectHoldGenerationCodeofCompleteArrangeNatureAndCanMaintainNature\n")
	prompt.WriteString("5. followGVAFrameworkofDevelopmentSpecificationAndbest practice\n")
	prompt.WriteString("6. ensureBeforeAfterEndFunctionCompleteArrangeToConnect\n")
	prompt.WriteString("7. AddmustRequireofErrorHandleAndDataVerify\n\n")
	prompt.WriteString("8. if neededvueRouteJump, PleaseUse menu_listergetCompleteArrangeRouteTable, AndAndRouteJumpUse router.push({\"name\":Frommenu_listerIngetofname})\n\n")
	prompt.WriteString("9. IfCurrentAllofvuePagecontentUnableSatisfyNeedRequest, ThenSelfRowbookWritevueFile, AndAndInvoke menu_creatorCreateMenuRecord\n\n")
	prompt.WriteString("10. if neededAPIInvoke, PleaseUse api_listergetapiTable, According toNeedRequestInvokeCorrespondingAPI\n\n")
	prompt.WriteString("11. IfCurrentAllAPIUnableSatisfyThenSelfRowbookWriteAPI, patchAllBeforeAfterEndGenerationCode, AndUse api_creatorCreateapiRecord\n\n")
	prompt.WriteString("12. NonetheoryBeforeAfterEndAllNotRequirefollowIntentiondeleteimportofcontent\n\n")
	prompt.WriteString("**PleaseBased onuser requirementAndAppearHaveFile, ProvideCompleteArrangeofGenerationCodeOptimizationPlan. **")

	return prompt.String()
}
