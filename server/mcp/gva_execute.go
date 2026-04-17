package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/mark3labs/mcp-go/mcp"
)

// register tool
func init() {
	RegisterTool(&GVAExecutor{})
}

// GVAExecutor GVACode Generator
type GVAExecutor struct{}

// ExecuteRequest ExecuteRequestStructureBody
type ExecuteRequest struct {
	ExecutionPlan ExecutionPlan `json:"executionPlan"` // ExecuteCountDraw
	Requirement   string        `json:"requirement"`   // OriginalNeedRequest(CanChoose, Used forLogRecord)
}

// ExecuteResponse ExecuteResponseStructureBody
type ExecuteResponse struct {
	Success        bool              `json:"success"`
	Message        string            `json:"message"`
	PackageID      uint              `json:"packageId,omitempty"`
	HistoryID      uint              `json:"historyId,omitempty"`
	Paths          map[string]string `json:"paths,omitempty"`
	GeneratedPaths []string          `json:"generatedPaths,omitempty"`
	NextActions    []string          `json:"nextActions,omitempty"`
}

// ExecutionPlan ExecuteCountDrawStructureBody
type ExecutionPlan struct {
	PackageName             string                            `json:"packageName"`
	PackageType             string                            `json:"packageType"` // "plugin" Or "package"
	NeedCreatedPackage      bool                              `json:"needCreatedPackage"`
	NeedCreatedModules      bool                              `json:"needCreatedModules"`
	NeedCreatedDictionaries bool                              `json:"needCreatedDictionaries"`
	PackageInfo             *request.SysAutoCodePackageCreate `json:"packageInfo,omitempty"`
	ModulesInfo             []*request.AutoCode               `json:"modulesInfo,omitempty"`
	Paths                   map[string]string                 `json:"paths,omitempty"`
	DictionariesInfo        []*DictionaryGenerateRequest      `json:"dictionariesInfo,omitempty"`
}

// New CreateGVAGenerationCodeGenerateExecutorUtility
func (g *GVAExecutor) New() mcp.Tool {
	return mcp.NewTool("gva_execute",
		mcp.WithDescription(`**GVAGenerationCodeGenerateExecutor:DirectExecuteGenerationCodeGenerate, No needsureAcknowledgeStepstep**

**Core features:**
According toNeedRequestPartAnalyzeAndCurrentofpackage infoJudgeYesNoInvoke, DirectGenerateGenerationCode. SupportBatchCreateMultipleModule, auto-createPackage, Module, Dictionaryetc.. 

**Use cases:**
Atgva_analyzegetDoneCurrentofpackage infoAndDictionaryInformationOfAfter, IfAlreadyViaPackageIncludeDoneCanByUseofPackageAndModule, thatThenNotRequireInvokeThismcp. According toPartAnalyzeResultDirectGenerateGenerationCode, SuitableUsed forautomation codeGenerateWorkflow. 

**Important notes:**
- WhenneedCreatedModules=truewhen true, module creation auto-generatesAPIAndMenu, NotShouldAgainInvokeapi_creatorAndmenu_creatorUtility
- FieldUsedictionary typeWhen, SystemWillAutomaticCheckAndcreate dictionary
- DictionaryCreateWillAtModuleCreateOfBeforeExecute
- WhenFieldconfigurationDonedataSourceAndassociation=2(one-to-manyAssociation)When, SystemWillAutomaticWillfieldTypeUpdateFor'array'`),
		mcp.WithObject("executionPlan",
			mcp.Description("ExecuteCountDraw, PackageIncludepackage info, ModuleAndDictionaryInformation"),
			mcp.Required(),
			mcp.Properties(map[string]interface{}{
				"packageName": map[string]interface{}{
					"type":        "string",
					"description": "package name(LowercasePrefix)",
				},
				"packageType": map[string]interface{}{
					"type":        "string",
					"description": "package Or plugin, IfUserRaiseToDoneUsePluginThenCreateplugin, IfUserNonespecialSetDescriptionThenOnelawChooseUsepackage",
					"enum":        []string{"package", "plugin"},
				},
				"needCreatedPackage": map[string]interface{}{
					"type":        "boolean",
					"description": "YesNoNeedCreatePackage, FortrueWhenpackageInfomustNeed",
				},
				"needCreatedModules": map[string]interface{}{
					"type":        "boolean",
					"description": "YesNoNeedcreate module, FortrueWhenmodulesInfomustNeed",
				},
				"needCreatedDictionaries": map[string]interface{}{
					"type":        "boolean",
					"description": "YesNoNeedcreate dictionary, FortrueWhendictionariesInfomustNeed",
				},
				"packageInfo": map[string]interface{}{
					"type":        "object",
					"description": "PackageCreation Info, WhenneedCreatedPackage=trueWhenmustNeed",
					"properties": map[string]interface{}{
						"desc":        map[string]interface{}{"type": "string", "description": "Packagedescription"},
						"label":       map[string]interface{}{"type": "string", "description": "display name"},
						"template":    map[string]interface{}{"type": "string", "description": "package Or plugin, IfUserRaiseToDoneUsePluginThenCreateplugin, IfUserNonespecialSetDescriptionThenOnelawChooseUsepackage", "enum": []string{"package", "plugin"}},
						"packageName": map[string]interface{}{"type": "string", "description": "package name"},
					},
				},
				"modulesInfo": map[string]interface{}{
					"type":        "array",
					"description": "ModuleconfigurationList, SupportBatchCreateMultipleModule",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"package":             map[string]interface{}{"type": "string", "description": "package name(LowercasePrefix, Example: userInfo)"},
							"tableName":           map[string]interface{}{"type": "string", "description": "Databasetable name(snakeShapeNameNamemethod,Example:user_info)"},
							"businessDB":          map[string]interface{}{"type": "string", "description": "business database(CankeepEmptyTableShowdefault)"},
							"structName":          map[string]interface{}{"type": "string", "description": "StructureBodyName(Largecamel caseExample:UserInfo)"},
							"packageName":         map[string]interface{}{"type": "string", "description": "file name"},
							"description":         map[string]interface{}{"type": "string", "description": "InTextdescription"},
							"abbreviation":        map[string]interface{}{"type": "string", "description": "alias"},
							"humpPackageName":     map[string]interface{}{"type": "string", "description": "file name(Smallcamel case), GeneralYesStructureBodyNameofSmallcamel caseExample:userInfo"},
							"gvaModel":            map[string]interface{}{"type": "boolean", "description": "YesNoUseGVAModel(FixedFortrue), AutomaticPackageIncludeID, CreatedAt, UpdatedAt, DeletedAtField"},
							"autoMigrate":         map[string]interface{}{"type": "boolean", "description": "YesNoAutomaticmigrateMoveDatabase"},
							"autoCreateResource":  map[string]interface{}{"type": "boolean", "description": "YesNoCreateResource(defaultForfalse)"},
							"autoCreateApiToSql":  map[string]interface{}{"type": "boolean", "description": "YesNoCreateAPI(defaultFortrue)"},
							"autoCreateMenuToSql": map[string]interface{}{"type": "boolean", "description": "YesNoCreateMenu(defaultFortrue)"},
							"autoCreateBtnAuth":   map[string]interface{}{"type": "boolean", "description": "YesNoCreateButton Permission(defaultForfalse)"},
							"onlyTemplate":        map[string]interface{}{"type": "boolean", "description": "YesNoOnlyTemplate(defaultForfalse)"},
							"isTree":              map[string]interface{}{"type": "boolean", "description": "whether tree structure(defaultForfalse)"},
							"treeJson":            map[string]interface{}{"type": "string", "description": "TreeShapeJSONField"},
							"isAdd":               map[string]interface{}{"type": "boolean", "description": "whether add operation(FixedForfalse)"},
							"generateWeb":         map[string]interface{}{"type": "boolean", "description": "YesNoGenerateFrontendGenerationCode"},
							"generateServer":      map[string]interface{}{"type": "boolean", "description": "YesNoGenerateAfterEndGenerationCode"},
							"fields": map[string]interface{}{
								"type":        "array",
								"description": "FieldList",
								"items": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"fieldName":       map[string]interface{}{"type": "string", "description": "field name(MustLargeWritePrefixExample:UserName)"},
										"fieldDesc":       map[string]interface{}{"type": "string", "description": "Fielddescription"},
										"fieldType":       map[string]interface{}{"type": "string", "description": "Fieldtype:string(String), richtext(richTextThis), int(ArrangeType), bool(booleanValue), float64(floatPointType), time.Time(Time), enum(Enum), picture(DocumentGraphSlice), pictures(MultipleGraphSlice), video(Viewfrequency), file(File), json(JSON), array(Array)"},
										"fieldJson":       map[string]interface{}{"type": "string", "description": "JSONTag,Example: userName"},
										"dataTypeLong":    map[string]interface{}{"type": "string", "description": "DataLength"},
										"comment":         map[string]interface{}{"type": "string", "description": "NoteExplain"},
										"columnName":      map[string]interface{}{"type": "string", "description": "DatabasecolumnName,Example: user_name"},
										"fieldSearchType": map[string]interface{}{"type": "string", "description": "Searchtype:=, !=, >, >=, <, <=, LIKE, BETWEEN, IN, NOT IN, NOT BETWEEN"},
										"fieldSearchHide": map[string]interface{}{"type": "boolean", "description": "is hiddenSearch"},
										"dictType":        map[string]interface{}{"type": "string", "description": "dictionary type, Usedictionary typeWhenSystemWillAutomaticCheckAndcreate dictionary"},
										"form":            map[string]interface{}{"type": "boolean", "description": "TableDocumentShow"},
										"table":           map[string]interface{}{"type": "boolean", "description": "TableformatShow"},
										"desc":            map[string]interface{}{"type": "boolean", "description": "DetailsShow"},
										"excel":           map[string]interface{}{"type": "boolean", "description": "importexport"},
										"require":         map[string]interface{}{"type": "boolean", "description": "YesNoRequired"},
										"defaultValue":    map[string]interface{}{"type": "string", "description": "defaultValue"},
										"errorText":       map[string]interface{}{"type": "string", "description": "ErrorPrompt"},
										"clearable":       map[string]interface{}{"type": "boolean", "description": "YesNoclearable"},
										"sort":            map[string]interface{}{"type": "boolean", "description": "YesNosort"},
										"primaryKey":      map[string]interface{}{"type": "boolean", "description": "YesNoPrimary Key(gvaModel=falseWhenMustHaveOnePieceFieldFortrue)"},
										"dataSource": map[string]interface{}{
											"type":        "object",
											"description": "data sourceconfiguration, Used forconfigurationFieldofrelated tableInformation. gettable namePrompt:CanAt server/model And plugin/xxx/model DirectoryDownviewCorrespondingModuleof TableName() Interface ImplementationgetActualtable name(Such As SysUser oftable nameFor sys_users). getDatabaseNamePrompt:MainDatabasePassCommonUse gva(defaultDatabaseIdentifier), MultipleDatabaseCanAt config.yaml of db-list configurationInviewCanUseDatabaseof alias-name Field, IfUserNotRaiseAndAssociationMultipleDatabaseInformationThenUsedefaultDatabase, defaultDatabaseIn this caseDown dbNameFillWriteEmpty",
											"properties": map[string]interface{}{
												"dbName":       map[string]interface{}{"type": "string", "description": "AssociationofDatabasename(defaultDatabasekeepEmpty)"},
												"table":        map[string]interface{}{"type": "string", "description": "Associationoftable name"},
												"label":        map[string]interface{}{"type": "string", "description": "Used forShowoffield name(Such Asname, titleetc.)"},
												"value":        map[string]interface{}{"type": "string", "description": "Used forSaveStoreofValuefield name(PassCommonYesid)"},
												"association":  map[string]interface{}{"type": "integer", "description": "associationtype:1=one-to-oneAssociation, 2=one-to-manyAssociation. one-to-oneAndone-to-manyofBeforePageofOneYesCurrentofEntity, IfHeOnlyAbleAssociationotherOnePieceEntityofOnePieceThenChooseUseone-to-one, IfHeNeedAssociationMultipleHeofAssociationEntityThenChooseUseone-to-many"},
												"hasDeletedAt": map[string]interface{}{"type": "boolean", "description": "related tableYesNoHavesoftdeleteField"},
											},
										},
										"checkDataSource": map[string]interface{}{"type": "boolean", "description": "YesNocheck data source, EnableAfterWillVerifyrelated tableofExistsNature"},
										"fieldIndexType":  map[string]interface{}{"type": "string", "description": "Indextype"},
									},
								},
							},
						},
					},
				},
				"paths": map[string]interface{}{
					"type":                 "object",
					"description":          "GenerateofFilepathMap",
					"additionalProperties": map[string]interface{}{"type": "string"},
				},
				"dictionariesInfo": map[string]interface{}{
					"type":        "array",
					"description": "DictionaryCreation Info, DictionaryCreateWillAtModuleCreateOfBeforeExecute",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"dictType":    map[string]interface{}{"type": "string", "description": "dictionary type, Used forIdentifierDictionaryofUniqueness"},
							"dictName":    map[string]interface{}{"type": "string", "description": "dictionary nameName, MustGenerate, DictionaryofChinese name"},
							"description": map[string]interface{}{"type": "string", "description": "dictionary description, DictionaryofUsageDescription"},
							"status":      map[string]interface{}{"type": "boolean", "description": "dictionary status:trueEnable, falseDisable"},
							"fieldDesc":   map[string]interface{}{"type": "string", "description": "Fielddescription, Used forAIReasonsolveFieldIncludemeaningAndGenerateAppropriateofChooseItem"},
							"options": map[string]interface{}{
								"type":        "array",
								"description": "DictionaryChooseItemList(CanChoose, if notProvideWillAccording tofieldDescauto-generatedefaultChooseItem)",
								"items": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"label": map[string]interface{}{"type": "string", "description": "Showname, UserSeeToofChooseItemName"},
										"value": map[string]interface{}{"type": "string", "description": "ChooseItemValue, ActualSaveStoreofValue"},
										"sort":  map[string]interface{}{"type": "integer", "description": "sortorder number; smaller sorts first"},
									},
								},
							},
						},
					},
				},
			}),
			mcp.AdditionalProperties(false),
		),
		mcp.WithString("requirement",
			mcp.Description("OriginalNeedRequestdescription(CanChoose, Used forLogRecord)"),
		),
	)
}

// Handle HandleExecuteRequest(RemovesureAcknowledgeStepstep)
func (g *GVAExecutor) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	executionPlanData, ok := request.GetArguments()["executionPlan"]
	if !ok {
		return nil, errors.New("Invalid parameters:executionPlan must be provided")
	}

	// parse execution plan
	planJSON, err := json.Marshal(executionPlanData)
	if err != nil {
		return nil, fmt.Errorf("parse execution planfailed: %v", err)
	}

	var plan ExecutionPlan
	err = json.Unmarshal(planJSON, &plan)
	if err != nil {
		return nil, fmt.Errorf("parse execution planfailed: %v\n\nPleaseensureExecutionPlanFormatCorrect, ReferenceUtilitydescriptionInStructureBodyFormatRequirement", err)
	}

	// validate execution plan completeness
	if err := g.validateExecutionPlan(&plan); err != nil {
		return nil, fmt.Errorf("ExecuteCountDrawVerifyfailed: %v", err)
	}

	// getOriginalNeedRequest(CanChoose)
	var originalRequirement string
	if reqData, ok := request.GetArguments()["requirement"]; ok {
		if reqStr, ok := reqData.(string); ok {
			originalRequirement = reqStr
		}
	}

	// DirectExecuteCreateOperation(NonesureAcknowledgeStepstep)
	result := g.executeCreation(ctx, &plan)

	// IfExecutesucceededAndHaveOriginalNeedRequest, ProvideGenerationCodecomplexCheckSuggest
	var reviewMessage string
	if result.Success && originalRequirement != "" {
		global.GVA_LOG.Info("ExecuteComplete, ReturnGenerateofFilepathprovideAIPerformGenerationCodecomplexCheck...")

		// BuildFilepathInformationprovideAIUse
		var pathsInfo []string
		for _, path := range result.GeneratedPaths {
			pathsInfo = append(pathsInfo, fmt.Sprintf("- %s", path))
		}

		reviewMessage = fmt.Sprintf("\n\n📁 AlreadyGenerateByDownFile:\n%s\n\n💡 Prompt:CanByCheckGenerateofGenerationCodeYesNoSatisfyOriginalNeedRequest. ", strings.Join(pathsInfo, "\n"))
	} else if originalRequirement == "" {
		reviewMessage = "\n\n💡 Prompt:Such AsNeedGenerationCodecomplexCheck, PleaseProvideOriginalNeedRequestdescription. "
	}

	// serialize response
	response := ExecuteResponse{
		Success:        result.Success,
		Message:        result.Message,
		PackageID:      result.PackageID,
		HistoryID:      result.HistoryID,
		Paths:          result.Paths,
		GeneratedPaths: result.GeneratedPaths,
		NextActions:    result.NextActions,
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serialize resultfailed: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("ExecuteResult:\n\n%s%s", string(responseJSON), reviewMessage)),
		},
	}, nil
}

// validateExecutionPlan validate execution plan completeness
func (g *GVAExecutor) validateExecutionPlan(plan *ExecutionPlan) error {
	if plan.PackageName == "" {
		return errors.New("packageName cannot be empty")
	}
	if plan.PackageType != "package" && plan.PackageType != "plugin" {
		return errors.New("packageType MustYes 'package' Or 'plugin'")
	}

	if plan.NeedCreatedPackage && plan.PackageInfo != nil && plan.PackageType != plan.PackageInfo.Template {
		return errors.New("packageType And packageInfo.template MustprotectHoldConsistent")
	}

	if plan.NeedCreatedPackage {
		if plan.PackageInfo == nil {
			return errors.New("When needCreatedPackage=true When, packageInfo cannot be empty")
		}
		if plan.PackageInfo.PackageName == "" {
			return errors.New("packageInfo.packageName cannot be empty")
		}
		if plan.PackageInfo.Template != "package" && plan.PackageInfo.Template != "plugin" {
			return errors.New("packageInfo.template MustYes 'package' Or 'plugin'")
		}
		if plan.PackageInfo.Label == "" {
			return errors.New("packageInfo.label cannot be empty")
		}
		if plan.PackageInfo.Desc == "" {
			return errors.New("packageInfo.desc cannot be empty")
		}
	}

	if plan.NeedCreatedModules {
		if len(plan.ModulesInfo) == 0 {
			return errors.New("When needCreatedModules=true When, modulesInfo cannot be empty")
		}

		for moduleIndex, moduleInfo := range plan.ModulesInfo {
			if moduleInfo.Package == "" {
				return fmt.Errorf("Module %d of package cannot be empty", moduleIndex+1)
			}
			if moduleInfo.StructName == "" {
				return fmt.Errorf("Module %d of structName cannot be empty", moduleIndex+1)
			}
			if moduleInfo.TableName == "" {
				return fmt.Errorf("Module %d of tableName cannot be empty", moduleIndex+1)
			}
			if moduleInfo.Description == "" {
				return fmt.Errorf("Module %d of description cannot be empty", moduleIndex+1)
			}
			if moduleInfo.Abbreviation == "" {
				return fmt.Errorf("Module %d of abbreviation cannot be empty", moduleIndex+1)
			}
			if moduleInfo.PackageName == "" {
				return fmt.Errorf("Module %d of packageName cannot be empty", moduleIndex+1)
			}
			if moduleInfo.HumpPackageName == "" {
				return fmt.Errorf("Module %d of humpPackageName cannot be empty", moduleIndex+1)
			}
			if len(moduleInfo.Fields) == 0 {
				return fmt.Errorf("Module %d of fields cannot be empty, At LeastNeedOnePieceField", moduleIndex+1)
			}

			for i, field := range moduleInfo.Fields {
				if field.FieldName == "" {
					return fmt.Errorf("Module %d Field %d of fieldName cannot be empty", moduleIndex+1, i+1)
				}
				if len(field.FieldName) > 0 {
					firstChar := string(field.FieldName[0])
					if firstChar >= "a" && firstChar <= "z" {
						moduleInfo.Fields[i].FieldName = strings.ToUpper(firstChar) + field.FieldName[1:]
					}
				}
				if field.FieldDesc == "" {
					return fmt.Errorf("Module %d Field %d of fieldDesc cannot be empty", moduleIndex+1, i+1)
				}
				if field.FieldType == "" {
					return fmt.Errorf("Module %d Field %d of fieldType cannot be empty", moduleIndex+1, i+1)
				}
				if field.FieldJson == "" {
					return fmt.Errorf("Module %d Field %d of fieldJson cannot be empty", moduleIndex+1, i+1)
				}
				if field.ColumnName == "" {
					return fmt.Errorf("Module %d Field %d of columnName cannot be empty", moduleIndex+1, i+1)
				}

				validFieldTypes := []string{"string", "int", "int64", "float64", "bool", "time.Time", "enum", "picture", "video", "file", "pictures", "array", "richtext", "json"}
				validType := false
				for _, validFieldType := range validFieldTypes {
					if field.FieldType == validFieldType {
						validType = true
						break
					}
				}
				if !validType {
					return fmt.Errorf("Module %d Field %d of fieldType '%s' not supported", moduleIndex+1, i+1, field.FieldType)
				}

				if field.FieldSearchType != "" {
					validSearchTypes := []string{"=", "!=", ">", ">=", "<", "<=", "LIKE", "BETWEEN", "IN", "NOT IN"}
					validSearchType := false
					for _, validSearchTypeValue := range validSearchTypes {
						if field.FieldSearchType == validSearchTypeValue {
							validSearchType = true
							break
						}
					}
					if !validSearchType {
						return fmt.Errorf("Module %d Field %d of fieldSearchType '%s' not supported", moduleIndex+1, i+1, field.FieldSearchType)
					}
				}

				if field.DataSource != nil {
					associationValue := field.DataSource.Association
					if associationValue == 2 && field.FieldType != "array" {
						global.GVA_LOG.Info(fmt.Sprintf("module %d field %d association=2, force fieldType to array", moduleIndex+1, i+1))
						moduleInfo.Fields[i].FieldType = "array"
					}
					if associationValue != 1 && associationValue != 2 {
						return fmt.Errorf("Module %d Field %d of dataSource.association MustYes 1 Or 2", moduleIndex+1, i+1)
					}
				}
			}

			if !moduleInfo.GvaModel {
				primaryKeyCount := 0
				for _, field := range moduleInfo.Fields {
					if field.PrimaryKey {
						primaryKeyCount++
					}
				}
				if primaryKeyCount == 0 {
					return fmt.Errorf("Module %d:When gvaModel=false When, MustHaveOnePieceFieldof primaryKey=true", moduleIndex+1)
				}
				if primaryKeyCount > 1 {
					return fmt.Errorf("Module %d:When gvaModel=false When, OnlyAbleHaveOnePieceFieldof primaryKey=true", moduleIndex+1)
				}
			} else {
				for i, field := range moduleInfo.Fields {
					if field.PrimaryKey {
						return fmt.Errorf("Module %d:When gvaModel=true When, Field %d of primaryKey ShouldThisFor false", moduleIndex+1, i+1)
					}
				}
			}
		}
	}

	return nil
}

// executeCreation ExecuteCreateOperation
func (g *GVAExecutor) executeCreation(ctx context.Context, plan *ExecutionPlan) *ExecuteResponse {
	result := &ExecuteResponse{
		Success:        false,
		Paths:          make(map[string]string),
		GeneratedPaths: []string{}, // initializeGenerateFilepathList
	}

	// NonetheorySuch AswhatAllFirstBuildDirectoryStructureInformation, ensurepathsalwaysReturn
	result.Paths = g.buildDirectoryStructure(plan)

	// RecordPrePeriodGenerateofFilepath
	result.GeneratedPaths = g.collectExpectedFilePaths(plan)

	if !plan.NeedCreatedModules {
		result.Success = true
		result.Message += "AlreadycolumnOutCurrentFunctionofInvolveofDirectoryStructureInformation; PleaseAtpathsInview; AndAndAtCorrespondingSpecifyFileimplemented inRelatedBusinessLogic; "
		return result
	}

	// CreatePackage(if needed)
	if plan.NeedCreatedPackage && plan.PackageInfo != nil {
		err := createAutoCodePackage(ctx, plan.PackageInfo)
		if err != nil {
			result.Message = fmt.Sprintf("CreatePackagefailed: %v", err)
			// That isMakeCreatePackagefailed, AlsoRequireReturnpathsInformation
			return result
		}
		result.Message += "PackageCreated successfully; "
	}

	// CreateSpecifyDictionary(if needed)
	if plan.NeedCreatedDictionaries && len(plan.DictionariesInfo) > 0 {
		dictResult := g.createDictionariesFromInfo(ctx, plan.DictionariesInfo)
		result.Message += dictResult
	}

	// Batchcreate dictionaryAndModule(if needed)
	if plan.NeedCreatedModules && len(plan.ModulesInfo) > 0 {
		// IterateAllModulePerformCreate
		for _, moduleInfo := range plan.ModulesInfo {

			// create module
			err := moduleInfo.Pretreatment()
			if err != nil {
				result.Message += fmt.Sprintf("Module %s Informationpreprocessingfailed: %v; ", moduleInfo.StructName, err)
				continue // continue with next module
			}

			err = createAutoCodeModule(ctx, *moduleInfo)
			if err != nil {
				result.Message += fmt.Sprintf("create module %s failed: %v; ", moduleInfo.StructName, err)
				continue // continue with next module
			}
			result.Message += fmt.Sprintf("Module %s Created successfully; ", moduleInfo.StructName)
		}

		result.Message += fmt.Sprintf("BatchCreateComplete, TotalHandle %d PieceModule; ", len(plan.ModulesInfo))

		// AddRepeatRequireReminder:NotRequireUseOthersMCPUtility
		result.Message += "\n\n⚠️ RepeatRequireReminder:\n"
		result.Message += "ModuleCreateFinished, APIAndMenuAlreadyauto-generate. PleaseNotRequireAgainInvokeByDownMCPUtility:\n"
		result.Message += "- api_creator:APIPermissionAlreadyAtModuleCreateWhenauto-generate\n"
		result.Message += "- menu_creator:FrontendMenuAlreadyAtModuleCreateWhenauto-generate\n"
		result.Message += "Such AsNeedUpdateAPIOrMenu, PleaseDirectAtSystemmanagementInterfaceInPerformconfiguration. \n"
	}

	result.Message += "AlreadyBuildDirectoryStructureInformation; "
	result.Success = true

	if result.Message == "" {
		result.Message = "ExecuteCountDrawComplete"
	}

	return result
}

// buildDirectoryStructure BuildDirectoryStructureInformation
func (g *GVAExecutor) buildDirectoryStructure(plan *ExecutionPlan) map[string]string {
	paths := make(map[string]string)

	// getconfigurationInformation
	autoCodeConfig := global.GVA_CONFIG.AutoCode

	// build basepath
	rootPath := autoCodeConfig.Root
	serverPath := autoCodeConfig.Server
	webPath := autoCodeConfig.Web
	moduleName := autoCodeConfig.Module

	// IfCountDrawInHavepackage name, UseCountDrawInpackage name, NoThenUsedefault
	packageName := "example"
	if plan.PackageName != "" {
		packageName = plan.PackageName
	}

	// IfCountDrawInHaveModuleInformation, getNo.OnePieceModuleofStructureNameAsdefaultValue
	structName := "ExampleStruct"
	if len(plan.ModulesInfo) > 0 && plan.ModulesInfo[0].StructName != "" {
		structName = plan.ModulesInfo[0].StructName
	}

	// According toPackagetypeBuildDifferentofpathStructure
	packageType := plan.PackageType
	if packageType == "" {
		packageType = "package" // defaultForpackageMode
	}

	// BuildServerpath
	if serverPath != "" {
		serverBasePath := fmt.Sprintf("%s/%s", rootPath, serverPath)

		if packageType == "plugin" {
			// Plugin Mode:AllFileAllAt /plugin/packageName/ in directory
			plugingBasePath := fmt.Sprintf("%s/plugin/%s", serverBasePath, packageName)

			// API path
			paths["api"] = fmt.Sprintf("%s/api", plugingBasePath)

			// Service path
			paths["service"] = fmt.Sprintf("%s/service", plugingBasePath)

			// Model path
			paths["model"] = fmt.Sprintf("%s/model", plugingBasePath)

			// Router path
			paths["router"] = fmt.Sprintf("%s/router", plugingBasePath)

			// Request path
			paths["request"] = fmt.Sprintf("%s/model/request", plugingBasePath)

			// Response path
			paths["response"] = fmt.Sprintf("%s/model/response", plugingBasePath)

			// Plugin specialHaveFile
			paths["plugin_main"] = fmt.Sprintf("%s/main.go", plugingBasePath)
			paths["plugin_config"] = fmt.Sprintf("%s/plugin.go", plugingBasePath)
			paths["plugin_initialize"] = fmt.Sprintf("%s/initialize", plugingBasePath)
		} else {
			// Package mode: classic directory layout
			// API path
			paths["api"] = fmt.Sprintf("%s/api/v1/%s", serverBasePath, packageName)

			// Service path
			paths["service"] = fmt.Sprintf("%s/service/%s", serverBasePath, packageName)

			// Model path
			paths["model"] = fmt.Sprintf("%s/model/%s", serverBasePath, packageName)

			// Router path
			paths["router"] = fmt.Sprintf("%s/router/%s", serverBasePath, packageName)

			// Request path
			paths["request"] = fmt.Sprintf("%s/model/%s/request", serverBasePath, packageName)

			// Response path
			paths["response"] = fmt.Sprintf("%s/model/%s/response", serverBasePath, packageName)
		}
	}

	// BuildFrontendpath(TwoTypeModeSame)
	if webPath != "" {
		webBasePath := fmt.Sprintf("%s/%s", rootPath, webPath)

		if packageType == "plugin" {
			// Plugin Mode:FrontendFileAlsoAt /plugin/packageName/ in directory
			pluginWebBasePath := fmt.Sprintf("%s/plugin/%s", webBasePath, packageName)

			// Vue Pagepath
			paths["vue_page"] = fmt.Sprintf("%s/view", pluginWebBasePath)

			// API path
			paths["vue_api"] = fmt.Sprintf("%s/api", pluginWebBasePath)
		} else {
			// Package mode: classic directory layout
			// Vue Pagepath
			paths["vue_page"] = fmt.Sprintf("%s/view/%s", webBasePath, packageName)

			// API path
			paths["vue_api"] = fmt.Sprintf("%s/api/%s", webBasePath, packageName)
		}
	}

	// AddModuleInformation
	paths["module"] = moduleName
	paths["package_name"] = packageName
	paths["package_type"] = packageType
	paths["struct_name"] = structName
	paths["root_path"] = rootPath
	paths["server_path"] = serverPath
	paths["web_path"] = webPath

	return paths
}

// collectExpectedFilePaths ReceivecollectPrePeriodGenerateofFilepath
func (g *GVAExecutor) collectExpectedFilePaths(plan *ExecutionPlan) []string {
	var paths []string

	// getDirectoryStructure
	dirPaths := g.buildDirectoryStructure(plan)

	// if neededcreate module, AddPrePeriodofFilepath
	if plan.NeedCreatedModules && len(plan.ModulesInfo) > 0 {
		for _, moduleInfo := range plan.ModulesInfo {
			structName := moduleInfo.StructName

			// AfterEndFile
			if apiPath, ok := dirPaths["api"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", apiPath, strings.ToLower(structName)))
			}
			if servicePath, ok := dirPaths["service"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", servicePath, strings.ToLower(structName)))
			}
			if modelPath, ok := dirPaths["model"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", modelPath, strings.ToLower(structName)))
			}
			if routerPath, ok := dirPaths["router"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", routerPath, strings.ToLower(structName)))
			}
			if requestPath, ok := dirPaths["request"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", requestPath, strings.ToLower(structName)))
			}
			if responsePath, ok := dirPaths["response"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.go", responsePath, strings.ToLower(structName)))
			}

			// FrontendFile
			if vuePage, ok := dirPaths["vue_page"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.vue", vuePage, strings.ToLower(structName)))
			}
			if vueApi, ok := dirPaths["vue_api"]; ok {
				paths = append(paths, fmt.Sprintf("%s/%s.js", vueApi, strings.ToLower(structName)))
			}
		}
	}

	return paths
}

// checkDictionaryExists check dictionaryYesNoExists
func (g *GVAExecutor) checkDictionaryExists(dictType string) (bool, error) {
	dictionary, err := findDictionaryByType(context.Background(), dictType)
	if err != nil {
		return false, err
	}
	return dictionary != nil, nil
}

// createDictionariesFromInfo According to DictionariesInfo create dictionary
func (g *GVAExecutor) createDictionariesFromInfo(ctx context.Context, dictionariesInfo []*DictionaryGenerateRequest) string {
	var messages []string

	messages = append(messages, fmt.Sprintf("StartCreate %d PieceSpecifyDictionary: ", len(dictionariesInfo)))

	for _, dictInfo := range dictionariesInfo {
		exists, err := g.checkDictionaryExists(dictInfo.DictType)
		if err != nil {
			messages = append(messages, fmt.Sprintf("check dictionary %s WhenOutError: %v; ", dictInfo.DictType, err))
			continue
		}

		if !exists {
			err = createDictionary(ctx, system.SysDictionary{
				Name:   dictInfo.DictName,
				Type:   dictInfo.DictType,
				Status: enabledBoolPointer(),
				Desc:   dictInfo.Description,
			})
			if err != nil {
				messages = append(messages, fmt.Sprintf("create dictionary %s failed: %v; ", dictInfo.DictType, err))
				continue
			}

			messages = append(messages, fmt.Sprintf("succeededcreate dictionary %s (%s); ", dictInfo.DictType, dictInfo.DictName))

			createdDict, err := findDictionaryByType(ctx, dictInfo.DictType)
			if err != nil {
				messages = append(messages, fmt.Sprintf("getcreated dictionaryfailed: %v; ", err))
				continue
			}
			if createdDict == nil {
				messages = append(messages, fmt.Sprintf("getcreated dictionaryfailed: %s; ", dictInfo.DictType))
				continue
			}

			if len(dictInfo.Options) > 0 {
				successCount := 0
				for _, option := range dictInfo.Options {
					dictionaryDetail := system.SysDictionaryDetail{
						Label:           option.Label,
						Value:           option.Value,
						Status:          enabledBoolPointer(),
						Sort:            option.Sort,
						SysDictionaryID: int(createdDict.ID),
					}

					err = createDictionaryDetail(ctx, dictionaryDetail)
					if err == nil {
						successCount++
					}
				}
				messages = append(messages, fmt.Sprintf("CreateDone %d PieceDictionaryChooseItem; ", successCount))
			}
		} else {
			messages = append(messages, fmt.Sprintf("Dictionary %s already exists; skip creation; ", dictInfo.DictType))
		}
	}

	return strings.Join(messages, "")
}
