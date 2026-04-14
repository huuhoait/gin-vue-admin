package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&DictionaryOptionsGenerator{})
}

type DictionaryOptionsGenerator struct{}

type DictionaryOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sort  int    `json:"sort"`
}

type DictionaryGenerateRequest struct {
	DictType    string             `json:"dictType"`
	FieldDesc   string             `json:"fieldDesc"`
	Options     []DictionaryOption `json:"options"`
	DictName    string             `json:"dictName"`
	Description string             `json:"description"`
}

type DictionaryGenerateResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	DictType     string `json:"dictType"`
	OptionsCount int    `json:"optionsCount"`
}

func (d *DictionaryOptionsGenerator) New() mcp.Tool {
	return mcp.NewTool("generate_dictionary_options",
		mcp.WithDescription("smartAbleGenerateDictionaryChooseItemAndAutomaticcreate dictionaryAndDictionaryDetails"),
		mcp.WithString("dictType",
			mcp.Required(),
			mcp.Description("dictionary type, Used forIdentifierDictionaryofUniqueness"),
		),
		mcp.WithString("fieldDesc",
			mcp.Required(),
			mcp.Description("Fielddescription, Used forAIReasonsolveFieldIncludemeaning"),
		),
		mcp.WithString("options",
			mcp.Required(),
			mcp.Description("DictionaryChooseItemJSONstring, format: [{\"label\":\"ShowName\",\"value\":\"Value\",\"sort\":1}]"),
		),
		mcp.WithString("dictName",
			mcp.Description("dictionary nameName, if notProvideWillauto-generate"),
		),
		mcp.WithString("description",
			mcp.Description("dictionary description"),
		),
	)
}

func (d *DictionaryOptionsGenerator) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dictType, ok := args["dictType"].(string)
	if !ok || dictType == "" {
		return nil, errors.New("dictType ParameterYesrequired")
	}
	fieldDesc, ok := args["fieldDesc"].(string)
	if !ok || fieldDesc == "" {
		return nil, errors.New("fieldDesc ParameterYesrequired")
	}
	optionsStr, ok := args["options"].(string)
	if !ok || optionsStr == "" {
		return nil, errors.New("options ParameterYesrequired")
	}

	var options []DictionaryOption
	if err := json.Unmarshal([]byte(optionsStr), &options); err != nil {
		return nil, fmt.Errorf("options invalid parameter format: %v", err)
	}
	if len(options) == 0 {
		return nil, errors.New("options cannot be empty")
	}

	req := &DictionaryGenerateRequest{
		DictType:    dictType,
		FieldDesc:   fieldDesc,
		Options:     options,
		DictName:    stringValue(args["dictName"]),
		Description: stringValue(args["description"]),
	}

	result, err := d.createDictionaryWithOptions(ctx, req)
	if err != nil {
		return nil, err
	}

	return textResultWithJSON("DictionaryChooseItemGenerateResult:", result)
}

func (d *DictionaryOptionsGenerator) createDictionaryWithOptions(ctx context.Context, req *DictionaryGenerateRequest) (*DictionaryGenerateResponse, error) {
	existingDict, err := findDictionaryByType(ctx, req.DictType)
	if err != nil {
		return nil, fmt.Errorf("check dictionaryYesNoExistsfailed: %v", err)
	}
	if existingDict != nil {
		return &DictionaryGenerateResponse{
			Success:      false,
			Message:      fmt.Sprintf("Dictionary %s already exists; skip creation", req.DictType),
			DictType:     req.DictType,
			OptionsCount: 0,
		}, nil
	}

	dictName := req.DictName
	if dictName == "" {
		dictName = d.generateDictionaryName(req.DictType, req.FieldDesc)
	}

	if err := createDictionary(ctx, system.SysDictionary{
		Name:   dictName,
		Type:   req.DictType,
		Status: enabledBoolPointer(),
		Desc:   req.Description,
	}); err != nil {
		return nil, fmt.Errorf("create dictionaryfailed: %v", err)
	}

	createdDict, err := findDictionaryByType(ctx, req.DictType)
	if err != nil {
		return nil, fmt.Errorf("getcreated dictionaryfailed: %v", err)
	}
	if createdDict == nil {
		return nil, fmt.Errorf("getcreated dictionaryfailed")
	}

	successCount := 0
	for _, option := range req.Options {
		err := createDictionaryDetail(ctx, system.SysDictionaryDetail{
			Label:           option.Label,
			Value:           option.Value,
			Status:          enabledBoolPointer(),
			Sort:            option.Sort,
			SysDictionaryID: int(createdDict.ID),
		})
		if err == nil {
			successCount++
		}
	}

	return &DictionaryGenerateResponse{
		Success:      true,
		Message:      fmt.Sprintf("succeededcreate dictionary %s, PackageInclude %d PieceChooseItem", req.DictType, successCount),
		DictType:     req.DictType,
		OptionsCount: successCount,
	}, nil
}

func (d *DictionaryOptionsGenerator) generateDictionaryName(dictType, fieldDesc string) string {
	if fieldDesc != "" {
		return fmt.Sprintf("%sDictionary", fieldDesc)
	}
	return fmt.Sprintf("%sDictionary", dictType)
}

func stringValue(value any) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
