package system

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common"
	system "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemResp "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
)

const (
	aiWorkflowMarkdownRootDir = "ai-workflow-docs"
	aiWorkflowAnalysisDir     = "analysis"
	aiWorkflowPromptDir       = "prompt-workflow"
)

func (s *aiWorkflowSession) DumpMarkdown(ctx context.Context, userID uint, info systemReq.SysAIWorkflowMarkdownDump) (result systemResp.AIWorkflowMarkdownDumpResult, err error) {
	if userID == 0 {
		return result, fmt.Errorf("user not logged in")
	}
	if info.Tab != "analysis" && info.Tab != "workflow" {
		return result, fmt.Errorf("unsupported session type")
	}

	root := strings.TrimSpace(global.GVA_CONFIG.AutoCode.Root)
	if root == "" {
		return result, fmt.Errorf("autocode.root is not configured")
	}

	session := system.SysAIWorkflowSession{
		GVA_MODEL:      global.GVA_MODEL{ID: info.ID},
		UserID:         userID,
		Tab:            info.Tab,
		Title:          strings.TrimSpace(info.Title),
		Summary:        strings.TrimSpace(info.Summary),
		ConversationID: strings.TrimSpace(info.ConversationID),
		MessageID:      strings.TrimSpace(info.MessageID),
		CurrentNodeID:  strings.TrimSpace(info.CurrentNodeID),
		Settings:       cloneJSONMap(info.Settings),
		FormData:       cloneJSONMap(info.FormData),
		ResultData:     cloneJSONMap(info.ResultData),
		Messages:       sanitizeMessages(info.Messages),
	}
	if strings.TrimSpace(session.Title) == "" {
		session.Title = s.titleFromMessages(session.Messages)
	}
	if strings.TrimSpace(session.Title) == "" {
		session.Title = s.titleFromForm(systemReq.SysAIWorkflowSessionUpsert{
			Tab:      info.Tab,
			FormData: info.FormData,
		})
	}
	if strings.TrimSpace(session.Summary) == "" {
		session.Summary = s.summaryFromResult(info.ResultData)
	}

	markdown := buildAIWorkflowMarkdown(session)
	if strings.TrimSpace(markdown) == "" {
		return result, fmt.Errorf("no content to dump")
	}

	targetDir := filepath.Join(root, aiWorkflowMarkdownRootDir, workflowMarkdownSubDir(info.Tab))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return result, err
	}

	fileName := buildWorkflowMarkdownFileName(info.Tab, session.Title, session.ID)
	filePath := filepath.Join(targetDir, fileName)
	if err := os.WriteFile(filePath, []byte(markdown), 0o644); err != nil {
		return result, err
	}

	relativePath, relErr := filepath.Rel(root, filePath)
	if relErr != nil {
		relativePath = filepath.Join(aiWorkflowMarkdownRootDir, workflowMarkdownSubDir(info.Tab), fileName)
	}

	return systemResp.AIWorkflowMarkdownDumpResult{
		FileName:     fileName,
		FilePath:     filePath,
		RelativePath: relativePath,
		Directory:    targetDir,
	}, nil
}

func workflowMarkdownSubDir(tab string) string {
	if tab == "analysis" {
		return aiWorkflowAnalysisDir
	}
	return aiWorkflowPromptDir
}

func buildWorkflowMarkdownFileName(tab, title string, sessionID uint) string {
	prefix := "prompt-workflow"
	if tab == "analysis" {
		prefix = "analysis"
	}

	stem := sanitizeWorkflowFileStem(title)
	if stem == "" {
		if sessionID > 0 {
			stem = fmt.Sprintf("session-%d", sessionID)
		} else {
			stem = "session"
		}
	}

	return fmt.Sprintf(
		"%s-%s-%s.md",
		time.Now().Format("20060102-150405"),
		prefix,
		stem,
	)
}

func sanitizeWorkflowFileStem(title string) string {
	var builder strings.Builder
	lastDash := false

	for _, r := range strings.TrimSpace(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(unicode.ToLower(r))
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if !lastDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

func buildAIWorkflowMarkdown(session system.SysAIWorkflowSession) string {
	if session.Tab == "analysis" {
		return buildAnalysisMarkdown(session)
	}
	return buildPromptWorkflowMarkdown(session)
}

func buildAnalysisMarkdown(session system.SysAIWorkflowSession) string {
	var builder strings.Builder

	writeMarkdownTitle(&builder, "# AI Requirements Analysis")
	writeMarkdownMeta(&builder,
		"Title", firstNonEmpty(session.Title, "Unnamed Requirement"),
		"Summary", firstNonEmpty(session.Summary, getString(session.ResultData, "summary")),
		"Session Type", "analysis",
		"Session ID", session.ConversationID,
		"Message ID", session.MessageID,
		"Node ID", session.CurrentNodeID,
	)

	writeMarkdownSection(&builder, "## Raw Input")
	writeMarkdownKeyValue(&builder,
		"Raw Requirement", getString(session.FormData, "requirement"),
		"Target Form", getString(session.FormData, "packageType"),
		"Business Scenario", getString(session.FormData, "businessScene"),
		"Extra Constraints", getString(session.FormData, "extraConstraints"),
		"Has Client Page", formatBool(getBool(session.FormData, "hasClientPage")),
		"Client Page Description", getString(session.FormData, "clientPageDescription"),
		"Client Extra Constraints", getString(session.FormData, "clientPageConstraints"),
	)

	writeMarkdownSection(&builder, "## Refined Requirements")
	writeMarkdownKeyValue(&builder,
		"Summary", getString(session.ResultData, "summary"),
		"Recommended Form", getString(session.ResultData, "recommendedPackageType"),
	)

	writeStringListSection(&builder, "### Pending Confirmations", getStringSlice(session.ResultData, "missingInfo"))
	writeStringListSection(&builder, "### Suggestions", getStringSlice(session.ResultData, "suggestions"))

	modules := getMapSlice(session.ResultData, "modules")
	if len(modules) > 0 {
		writeMarkdownSection(&builder, "### Module Breakdown")
		for index, module := range modules {
			builder.WriteString(fmt.Sprintf("#### %d. %s\n\n", index+1, firstNonEmpty(getString(module, "label"), getString(module, "name"), fmt.Sprintf("Module %d", index+1))))
			writeMarkdownKeyValue(&builder,
				"Module ID", getString(module, "name"),
				"Module Description", getString(module, "description"),
			)

			fields := getMapSlice(module, "fields")
			if len(fields) > 0 {
				builder.WriteString("| Field | ID | Type | Required | Description |\n")
				builder.WriteString("| --- | --- | --- | --- | --- |\n")
				for _, field := range fields {
					builder.WriteString(fmt.Sprintf(
						"| %s | %s | %s | %s | %s |\n",
						markdownCell(firstNonEmpty(getString(field, "label"), getString(field, "name"), "-")),
						markdownCell(firstNonEmpty(getString(field, "name"), "-")),
						markdownCell(firstNonEmpty(getString(field, "type"), "string")),
						markdownCell(formatBool(getBool(field, "required"))),
						markdownCell(firstNonEmpty(getString(field, "description"), "-")),
					))
				}
				builder.WriteString("\n")
			}
		}
	}

	clientPages := getMapSlice(session.ResultData, "clientPages")
	if len(clientPages) > 0 {
		writeMarkdownSection(&builder, "### Client Pages")
		for index, page := range clientPages {
			builder.WriteString(fmt.Sprintf("#### %d. %s\n\n", index+1, firstNonEmpty(getString(page, "label"), getString(page, "name"), fmt.Sprintf("Page %d", index+1))))
			writeMarkdownKeyValue(&builder,
				"Page ID", getString(page, "name"),
				"Page Type", getString(page, "pageType"),
				"Page Description", getString(page, "description"),
			)
			writeStringListSection(&builder, "Target Modules", getStringSlice(page, "targetModules"))
			writeStringListSection(&builder, "Interactions", getStringSlice(page, "interactions"))
			writeStringListSection(&builder, "Field Relations", getStringSlice(page, "relations"))
		}
	}

	writeMarkdownAppendix(&builder, session.ResultData)
	return strings.TrimSpace(builder.String()) + "\n"
}

func buildPromptWorkflowMarkdown(session system.SysAIWorkflowSession) string {
	var builder strings.Builder

	writeMarkdownTitle(&builder, "# Prompt Workflow")
	writeMarkdownMeta(&builder,
		"Title", firstNonEmpty(session.Title, "Unnamed Workflow"),
		"Summary", firstNonEmpty(session.Summary, getString(session.ResultData, "summary")),
		"Session Type", "workflow",
		"Session ID", session.ConversationID,
		"Message ID", session.MessageID,
		"Node ID", session.CurrentNodeID,
	)

	writeMarkdownSection(&builder, "## Input Context")
	writeMarkdownKeyValue(&builder,
		"Source Requirement", getString(session.FormData, "source"),
		"Workflow Type", getString(session.FormData, "flowType"),
		"Extra Constraints", getString(session.FormData, "extraConstraints"),
	)

	writeMarkdownSection(&builder, "## Prompt Workflow")
	writeMarkdownKeyValue(&builder, "Summary", getString(session.ResultData, "summary"))

	steps := getMapSlice(session.ResultData, "steps")
	if len(steps) == 0 {
		rawText := getString(session.ResultData, "rawText")
		if strings.TrimSpace(rawText) != "" {
			builder.WriteString(rawText)
			builder.WriteString("\n\n")
		}
	} else {
		for index, step := range steps {
			builder.WriteString(fmt.Sprintf("### %d. %s\n\n", index+1, firstNonEmpty(getString(step, "title"), fmt.Sprintf("Step %d", index+1))))
			writeMarkdownKeyValue(&builder,
				"Goal", getString(step, "goal"),
				"Suggested Tool", getString(step, "suggestedTool"),
				"Auto-executable", formatBool(getBool(step, "autoExecutable")),
				"Expected Output", getString(step, "expectedOutput"),
			)
			prompt := getString(step, "prompt")
			if strings.TrimSpace(prompt) != "" {
				builder.WriteString("#### Prompt\n\n")
				builder.WriteString("```text\n")
				builder.WriteString(prompt)
				builder.WriteString("\n```\n\n")
			}
		}
	}

	writeMarkdownAppendix(&builder, session.ResultData)
	return strings.TrimSpace(builder.String()) + "\n"
}

func writeMarkdownTitle(builder *strings.Builder, title string) {
	builder.WriteString(title)
	builder.WriteString("\n\n")
}

func writeMarkdownMeta(builder *strings.Builder, pairs ...string) {
	for i := 0; i+1 < len(pairs); i += 2 {
		if strings.TrimSpace(pairs[i+1]) == "" {
			continue
		}
		builder.WriteString(fmt.Sprintf("- **%s**: %s\n", pairs[i], pairs[i+1]))
	}
	builder.WriteString("\n")
}

func writeMarkdownSection(builder *strings.Builder, title string) {
	builder.WriteString(title)
	builder.WriteString("\n\n")
}

func writeMarkdownKeyValue(builder *strings.Builder, pairs ...string) {
	hasValue := false
	for i := 0; i+1 < len(pairs); i += 2 {
		if strings.TrimSpace(pairs[i+1]) == "" {
			continue
		}
		hasValue = true
		builder.WriteString(fmt.Sprintf("- **%s**: %s\n", pairs[i], pairs[i+1]))
	}
	if hasValue {
		builder.WriteString("\n")
	}
}

func writeStringListSection(builder *strings.Builder, title string, list []string) {
	if len(list) == 0 {
		return
	}
	builder.WriteString(title)
	builder.WriteString("\n\n")
	for _, item := range list {
		if strings.TrimSpace(item) == "" {
			continue
		}
		builder.WriteString("- ")
		builder.WriteString(item)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")
}

func writeMarkdownAppendix(builder *strings.Builder, resultData common.JSONMap) {
	rawText := getString(resultData, "rawText")
	rawJSON := getString(resultData, "rawJson")

	if strings.TrimSpace(rawText) != "" {
		builder.WriteString("## Raw Text\n\n")
		builder.WriteString("```text\n")
		builder.WriteString(rawText)
		builder.WriteString("\n```\n\n")
	}

	if strings.TrimSpace(rawJSON) != "" {
		builder.WriteString("## Raw Structured Data\n\n")
		builder.WriteString("```json\n")
		builder.WriteString(rawJSON)
		builder.WriteString("\n```\n")
	}
}

func getString(data map[string]interface{}, key string) string {
	if len(data) == 0 {
		return ""
	}
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", value))
	}
}

func getBool(data map[string]interface{}, key string) bool {
	if len(data) == 0 {
		return false
	}
	value, ok := data[key]
	if !ok || value == nil {
		return false
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(strings.TrimSpace(typed), "true")
	default:
		return false
	}
}

func getStringSlice(data map[string]interface{}, key string) []string {
	value, ok := data[key]
	if !ok || value == nil {
		return nil
	}
	return toStringSlice(value)
}

func getMapSlice(data map[string]interface{}, key string) []map[string]interface{} {
	value, ok := data[key]
	if !ok || value == nil {
		return nil
	}
	return toMapSlice(value)
}

func toStringSlice(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if strings.TrimSpace(item) != "" {
				result = append(result, strings.TrimSpace(item))
			}
		}
		return result
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			text := strings.TrimSpace(fmt.Sprintf("%v", item))
			if text != "" {
				result = append(result, text)
			}
		}
		return result
	default:
		text := strings.TrimSpace(fmt.Sprintf("%v", value))
		if text == "" || text == "<nil>" {
			return nil
		}
		return []string{text}
	}
}

func toMapSlice(value interface{}) []map[string]interface{} {
	switch typed := value.(type) {
	case []map[string]interface{}:
		return typed
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			if row, ok := item.(map[string]interface{}); ok {
				result = append(result, row)
			}
		}
		return result
	default:
		return nil
	}
}

func formatBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func markdownCell(value string) string {
	text := strings.TrimSpace(value)
	if text == "" {
		return "-"
	}
	text = strings.ReplaceAll(text, "\n", "<br>")
	text = strings.ReplaceAll(text, "|", "\\|")
	return text
}
