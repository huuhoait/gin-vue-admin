package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&RequirementAnalyzer{})
}

type RequirementAnalyzer struct{}

// RequirementAnalysisRequest NeedRequestPartAnalyzeRequest
type RequirementAnalysisRequest struct {
	UserRequirement string `json:"userRequirement"`
}

// RequirementAnalysisResponse NeedRequestPartAnalyzeResponse
type RequirementAnalysisResponse struct {
	AIPrompt string `json:"aiPrompt"` // ToAIofPromptWord
}

// New returns tool registration information
func (t *RequirementAnalyzer) New() mcp.Tool {
	return mcp.NewTool("requirement_analyzer",
		mcp.WithDescription(`** smartAbleNeedRequestPartAnalyzeAndModuleDesignUtility - FirstChooseInPortUtility(HighestPriority)**

** Important notes:ThisYesAllMCPUtilityofFirstChooseInPort, PleaseExcellentFirstUse!**

** CoreCapability:**
AsfundDeepSystemarchitectureTeacher, smartAbleanalyze user requirementAndAutomaticDesignCompleteArrange modulearchitecture

** Core features:**
1. **smartAbleNeedRequestdestructure**:Deepdegreeanalyze user requirement, IdentifyCoreBusinessEntity, BusinessWorkflow, DataRelationship
2. **AutomaticModuleDesign**:Based onNeedRequestPartAnalyze, smartAblesureSetNeedHow ManyPieceModuleAndEachModuleFunction
3. **FieldsmartAblepushGuide**:ForEachModuleAutomaticDesignDetailsField, PackageIncludeDatatype, association, DictionaryNeedRequest
4. **architectureOptimizationSuggest**:ProvideModulesplitPart, AssociationDesign, ExtensionNatureetc.SpecialIndustrySuggest

** InputOutcontent:**
- ModuleQuantityAndarchitectureDesign
- EachModuleofDetailsFieldClearDocument
- DatatypeAndassociationDesign
- DictionaryNeedRequestAndtypeDefine
- ModuleIntervalRelationshipGraphAndExtensionSuggest

** SuitableUseScenario:**
- user requirementdescriptionNotCompleteArrange, NeedsmartAblepatchAll
- ComplexBusinessSystem modulearchitectureDesign
- NeedSpecialIndustryofDatabaseDesignSuggest
- ThinkRequirerapidly buildProduceLevelBusinessSystem

** recommendWorkflow:**
 requirement_analyzer → gva_analyze → gva_execute → OthersAssisthelpUtility
 
 `),
		mcp.WithString("userRequirement",
			mcp.Required(),
			mcp.Description("UserofNeedRequestdescription, SupportSelfThenLanguage, Such As:'IRequireDoOnePiececatterymanagementSystem, UseComerecordIncatofInformation, AndAndRecordEveryOnlycatEveryDayofActivityInformation'"),
		),
	)
}

// Handle HandleUtilityInvoke
func (t *RequirementAnalyzer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userRequirement, ok := request.GetArguments()["userRequirement"].(string)
	if !ok || userRequirement == "" {
		return nil, errors.New("Invalid parameters:userRequirement MustYesnon-empty string")
	}

	// analyze user requirement
	analysisResponse, err := t.analyzeRequirement(userRequirement)
	if err != nil {
		return nil, fmt.Errorf("NeedRequestPartAnalyzefailed: %v", err)
	}

	// serialize response
	responseData, err := json.Marshal(analysisResponse)
	if err != nil {
		return nil, fmt.Errorf("serialize responsefailed: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(responseData)),
		},
	}, nil
}

// analyzeRequirement analyze user requirement - SpecialNoteAtAINeedRequestPass
func (t *RequirementAnalyzer) analyzeRequirement(userRequirement string) (*RequirementAnalysisResponse, error) {
	// GenerateAIPromptWord - ThisYesUniqueFunction
	aiPrompt := t.generateAIPrompt(userRequirement)

	return &RequirementAnalysisResponse{
		AIPrompt: aiPrompt,
	}, nil
}

// generateAIPrompt GenerateAIPromptWord - smartAblePartAnalyzeNeedRequestAndsureSetModuleStructure
func (t *RequirementAnalyzer) generateAIPrompt(userRequirement string) string {
	prompt := fmt.Sprintf(`# smartAbleNeedRequestPartAnalyzeAndModuleDesignTask

## UserOriginalNeedRequest
%s

## CoreTask
YouNeedAsOnePiecefundDeepofSystemarchitectureTeacher, Deepdegreeanalyze user requirement, smartAbleDesignOutCompleteArrange modulearchitecture. 

## PartAnalyzeStepstep

### Step One:NeedRequestdestructurePartAnalyze
Pleasecarefulanalyze user requirement, IdentifyOut:
1. **CoreBusinessEntity**(Such As:User, Product, Order, vaccine, petItemetc.)
2. **BusinessWorkflow**(Such As:Register, PurchaseBuy, Record, managementetc.)
3. **DataRelationship**(EntityIntervalofassociation)
4. **FunctionModule**(NeedWhichIndependentofmanagementModule)

### Step Two:ModulearchitectureDesign
Based onNeedRequestPartAnalyze, DesignOutModulearchitecture, FormatSuch AsDown:

**Module1:[Modulename]**
- Functiondescription:[ThisModuleofCoreFunction]
- MainRequireField:[columnOutkeywordSegment, NotebrightDatatype]
- association:[AndOthersModuleofRelationship, clearone-to-one/one-to-many]
- DictionaryNeedRequest:[NeedWhichdictionary type]

**Module2:[Modulename]**
- Functiondescription:[ThisModuleofCoreFunction]
- MainRequireField:[columnOutkeywordSegment, NotebrightDatatype]
- association:[AndOthersModuleofRelationship]
- DictionaryNeedRequest:[NeedWhichdictionary type]

**...**

### No.ThreeStep:FieldDetailsDesign
ForEachModuleDetailsDesignField:

#### Module1FieldClearDocument:
- field name1 (Datatype) - Fielddescription [YesNoRequired] [association info/dictionary type]
- field name2 (Datatype) - Fielddescription [YesNoRequired] [association info/dictionary type]
- ...

#### Module2FieldClearDocument:
- field name1 (Datatype) - Fielddescription [YesNoRequired] [association info/dictionary type]
- ...

## smartAblePartAnalyzeInstructionGuideOriginalThen

### ModulesplitPartOriginalThen
1. **DocumentOneOccupationresponsibility**:EachModuleOnlyResponsibleOnePieceCoreBusinessEntity
2. **DataCompleteArrangeNature**:RelatedDataShouldThisAtSameOneModuleIn
3. **BusinessIndependentNature**:ModuleShouldThisAble ToIndependentCompletespecialSetBusinessFunction
4. **ExtensionNatureConsider**:ForNotComeFunctionExtensionReserveEmptyInterval

### FieldDesignOriginalThen
1. **mustRequireNature**:OnlyPackageIncludeBusinessrequiredField
2. **SpecificationNature**:followDatabaseDesignSpecification
3. **AssociationNature**:CorrectIdentifyEntityIntervalRelationship
4. **Dictionarytransform**:status, typeetc.EnumValueUseDictionary

### associationIdentify
- **one-to-one**:OnePieceEntityOnlyAbleAssociationotherOnePieceEntityofOnePieceRecord
- **one-to-many**:OnePieceEntityCanByAssociationotherOnePieceEntityofMultipleRecord
- **MultipleToMultiple**:ApprovedInIntervalTableImplementComplexAssociation

## SpecialScenarioHandle

### ComplexEntityIdentify
WhenUserRaiseToCertainPiececonceptWhen, RequireJudgeItYesNoNeedIndependentModule:
- **DictionaryHandle**:SimpleOf the orderCommonSeeofstatus, type(Such As:Toggle, Gender, Completestatusetc.)
- **IndependentModule**:ComplexEntity(Such As:vaccinemanagement, petItemfileCase, NoteshootRecord)

## InputOutRequirement

### MustPackageIncludeofInformation
1. **ModuleQuantity**:clearNeedFewPieceModule
2. **ModuleRelationshipGraph**:UseTextCharacterdescriptionModuleIntervalRelationship
3. **CoreField**:EachModuleofkeywordSegment(At Least5-10Piece)
4. **Datatype**:string, int, bool, time.Time, float64etc.
5. **AssociationDesign**:clearWhichFieldYesAssociationField
6. **DictionaryNeedRequest**:columnOutNeedCreateofdictionary type

### strictly followUserInputIn
- IfUserProvideDoneToolBodyField, **MustUse**UserProvideofField
- IfUserProvideDoneSQLFile, **strictAccording to**SQLStructureDesign
- **NotRequire**followIntentiondiverge, **NotRequire**AddUserNotRaiseAndofFunction
---

**AppearAtPleaseStartDeepdegreeanalyze user requirement:"%s"**

PleaseAccording toAboveFrameworkPerformSystemNaturePartAnalyze, ensureInputOut moduleDesignalreadySatisfyCurrentNeedRequest, AgainToolPreparegoodGoodofExtensionNature. `, userRequirement, userRequirement)

	return prompt
}
