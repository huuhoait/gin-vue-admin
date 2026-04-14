package system

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

func (autoApi *AutoCodeApi) LLMAutoSSE(c *gin.Context) {
	var llm common.JSONMap
	if err := c.ShouldBindJSON(&llm); err != nil {
		global.GVA_LOG.Error("LLMAutoSSE parameter binding failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if llm == nil {
		llm = common.JSONMap{}
	}
	llm["response_mode"] = "streaming"
	global.GVA_LOG.Info("LLMAutoSSE received request", zap.Any("mode", llm["mode"]))

	if err := autoApi.streamLLMAsSSE(c, llm); err != nil {
		global.GVA_LOG.Error("LLM SSE proxy failed!", zap.Error(err))
		if c.Writer.Written() {
			writeLLMStreamError(c, err)
			return
		}
		response.FailWithMessage(err.Error(), c)
	}
}

func (autoApi *AutoCodeApi) streamLLMAsSSE(c *gin.Context, llm common.JSONMap) error {
	res, err := autoCodeService.LLMAutoStream(c.Request.Context(), llm)
	if err != nil {
		return fmt.Errorf("failed to call upstream LLM: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("upstream LLM streaming service returned non-2xx: status=%d content-type=%s read-body-err=%w", res.StatusCode, res.Header.Get("Content-Type"), readErr)
		}
		return fmt.Errorf("upstream LLM streaming service returned non-2xx: status=%d content-type=%s body=%s", res.StatusCode, res.Header.Get("Content-Type"), previewResponseBody(body))
	}

	ct := res.Header.Get("Content-Type")
	global.GVA_LOG.Info("LLMAutoSSE upstream returned successfully, starting SSE streaming relay",
		zap.Int("status", res.StatusCode),
		zap.String("content-type", ct))

	// If upstream does not return an SSE stream (possibly blocking mode JSON), read and relay directly
	if !strings.Contains(ct, "text/event-stream") && !strings.Contains(ct, "text/plain") {
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("failed to read upstream non-streaming response: %w", readErr)
		}
		global.GVA_LOG.Warn("LLMAutoSSE upstream returned non-SSE stream, Content-Type: "+ct+", will relay as single event",
			zap.String("body_preview", previewResponseBody(body)))

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			return errors.New("current response does not support streaming output")
		}
		prepareSSEHeaders(c)
		c.Status(http.StatusOK)

		var payload any
		if err := json.Unmarshal(body, &payload); err != nil {
			payload = string(body)
		}
		if err := renderSSE(c, sse.Event{Event: "message", Data: payload}); err != nil {
			return err
		}
		if err := renderSSE(c, sse.Event{Event: "done", Data: gin.H{"done": true}}); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return errors.New("current response does not support streaming output")
	}

	prepareSSEHeaders(c)
	c.Status(http.StatusOK)
	flusher.Flush()

	reader := bufio.NewReader(res.Body)
	lines := make([]string, 0, 8)
	blockCount := 0

	global.GVA_LOG.Info("LLMAutoSSE starting to read upstream stream data...")

	for {
		global.GVA_LOG.Debug("LLMAutoSSE waiting to read next line...")
		line, readErr := reader.ReadString('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			global.GVA_LOG.Error("LLMAutoSSE failed to read upstream stream", zap.Int("blocks_relayed", blockCount), zap.Error(readErr))
			return fmt.Errorf("failed to read upstream streaming response: %w", readErr)
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if len(lines) > 0 {
				blockCount++
				if blockCount <= 3 {
					global.GVA_LOG.Debug("LLMAutoSSE relaying SSE block", zap.Int("block", blockCount), zap.Strings("lines", lines))
				}
			}
			if err := emitSSEBlock(c, lines); err != nil {
				return err
			}
			lines = lines[:0]
		} else {
			lines = append(lines, line)
		}

		if errors.Is(readErr, io.EOF) {
			if err := emitSSEBlock(c, lines); err != nil {
				return err
			}
			if err := renderSSE(c, sse.Event{
				Event: "done",
				Data:  gin.H{"done": true},
			}); err != nil {
				return err
			}
			flusher.Flush()
			global.GVA_LOG.Info("LLMAutoSSE streaming relay completed", zap.Int("total_blocks", blockCount))
			return nil
		}
	}
}

func prepareSSEHeaders(c *gin.Context) {
	header := c.Writer.Header()
	header.Set("Content-Type", "text/event-stream; charset=utf-8")
	header.Set("Cache-Control", "no-cache, no-transform")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
}

func emitSSEBlock(c *gin.Context, lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	eventName := "message"
	eventID := ""
	dataLines := make([]string, 0, len(lines))

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "event:"):
			eventName = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "id:"):
			eventID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		case strings.HasPrefix(line, "data:"):
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
		}
	}

	rawData := strings.TrimSpace(strings.Join(dataLines, "\n"))
	if rawData == "" {
		return nil
	}
	if rawData == "[DONE]" {
		return renderSSE(c, sse.Event{
			Id:    eventID,
			Event: "done",
			Data:  gin.H{"done": true},
		})
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(rawData), &payload); err != nil {
		payload = rawData
	}

	return renderSSE(c, sse.Event{
		Id:    eventID,
		Event: eventName,
		Data:  payload,
	})
}

func renderSSE(c *gin.Context, event sse.Event) error {
	if err := event.Render(c.Writer); err != nil {
		return fmt.Errorf("failed to write SSE event: %w", err)
	}
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
