package system

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	commonResp "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"github.com/goccy/go-json"
)

func (s *AutoCodeService) LLMAuto(ctx context.Context, llm common.JSONMap) (interface{}, error) {
	path, err := buildLLMAutoPath(llm)
	if err != nil {
		return nil, err
	}

	res, err := request.HttpRequestWithContextAndTimeout(
		ctx,
		path,
		http.MethodPost,
		nil,
		nil,
		llm,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call upstream LLM service: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read LLM response: %w", err)
	}

	bodyPreview := previewResponseBody(body)
	contentType := res.Header.Get("Content-Type")
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("upstream LLM service returned non-2xx: status=%d content-type=%s body=%s", res.StatusCode, contentType, bodyPreview)
	}

	var resStruct commonResp.Response
	if err = json.Unmarshal(body, &resStruct); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: status=%d content-type=%s body=%s err=%w", res.StatusCode, contentType, bodyPreview, err)
	}

	if resStruct.Code != commonResp.SUCCESS {
		return nil, fmt.Errorf("LLM service returned business error: code=%d msg=%s body=%s", resStruct.Code, resStruct.Msg, bodyPreview)
	}

	return resStruct.Data, nil
}

func (s *AutoCodeService) LLMAutoStream(ctx context.Context, llm common.JSONMap) (*http.Response, error) {
	path, err := buildLLMAutoPath(llm)
	if err != nil {
		return nil, err
	}

	payload := cloneLLMAutoJSONMap(llm)
	responseMode := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", payload["response_mode"])))
	if responseMode == "" {
		payload["response_mode"] = "streaming"
	}

	res, err := request.HttpRequestWithContextAndTimeout(
		ctx,
		path,
		http.MethodPost,
		map[string]string{
			"Accept":          "text/event-stream",
			"Accept-Encoding": "identity", // Disable gzip to prevent SSE stream compression causing buffer stalls
			"Cache-Control":   "no-cache",
		},
		nil,
		payload,
		-1, // Do not set client.Timeout; SSE stream lifecycle is controlled by ctx
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call upstream LLM streaming service: %w", err)
	}
	return res, nil
}

func buildLLMAutoPath(llm common.JSONMap) (string, error) {
	if global.GVA_CONFIG.AutoCode.AiPath == "" {
		return "", errors.New("please go to the plugin marketplace personal center to obtain AiPath and fill it in config.yaml")
	}

	mode := strings.TrimSpace(fmt.Sprintf("%v", llm["mode"]))
	if mode == "" {
		return "", errors.New("llmAuto missing mode parameter")
	}

	return strings.ReplaceAll(global.GVA_CONFIG.AutoCode.AiPath, "{FUNC}", mode), nil
}

func cloneLLMAutoJSONMap(src common.JSONMap) common.JSONMap {
	dst := make(common.JSONMap, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func previewResponseBody(body []byte) string {
	text := strings.TrimSpace(string(body))
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return "<empty>"
	}
	runes := []rune(text)
	if len(runes) > 300 {
		return string(runes[:300]) + "..."
	}
	return text
}
