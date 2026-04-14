package system

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type AutoCodeApi struct{}

func (autoApi *AutoCodeApi) GetDB(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbs, err := autoCodeService.Database(businessDB).GetDB(businessDB)
	var dbList []map[string]interface{}
	for _, db := range global.GVA_CONFIG.DBList {
		item := map[string]interface{}{
			"aliasName": db.AliasName,
			"dbName":    db.Dbname,
			"disable":   db.Disable,
			"dbtype":    db.Type,
		}
		dbList = append(dbList, item)
	}
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"dbs": dbs, "dbList": dbList}, "Retrieved successfully", c)
}

func (autoApi *AutoCodeApi) GetTables(c *gin.Context) {
	dbName := c.Query("dbName")
	businessDB := c.Query("businessDB")
	if dbName == "" {
		dbName = *global.GVA_ACTIVE_DBNAME
		if businessDB != "" {
			for _, db := range global.GVA_CONFIG.DBList {
				if db.AliasName == businessDB {
					dbName = db.Dbname
				}
			}
		}
	}

	tables, err := autoCodeService.Database(businessDB).GetTables(businessDB, dbName)
	if err != nil {
		global.GVA_LOG.Error("Failed to query tables!", zap.Error(err))
		response.FailWithMessage("Failed to query tables", c)
		return
	}
	response.OkWithDetailed(gin.H{"tables": tables}, "Retrieved successfully", c)
}

func (autoApi *AutoCodeApi) GetColumn(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbName := c.Query("dbName")
	if dbName == "" {
		dbName = *global.GVA_ACTIVE_DBNAME
		if businessDB != "" {
			for _, db := range global.GVA_CONFIG.DBList {
				if db.AliasName == businessDB {
					dbName = db.Dbname
				}
			}
		}
	}
	tableName := c.Query("tableName")
	columns, err := autoCodeService.Database(businessDB).GetColumn(businessDB, tableName, dbName)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"columns": columns}, "Retrieved successfully", c)
}

func (autoApi *AutoCodeApi) LLMAuto(c *gin.Context) {
	var llm common.JSONMap
	if err := c.ShouldBindJSON(&llm); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if shouldStreamLLM(c, llm) {
		if err := autoApi.proxyLLMStream(c, llm); err != nil {
			global.GVA_LOG.Error("LLM streaming proxy failed!", zap.Error(err))
			if c.Writer.Written() {
				writeLLMStreamError(c, err)
				return
			}
			response.FailWithMessage(err.Error(), c)
		}
		return
	}

	data, err := autoCodeService.LLMAuto(c.Request.Context(), llm)
	if err != nil {
		global.GVA_LOG.Error("LLM generation failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func shouldStreamLLM(c *gin.Context, llm common.JSONMap) bool {
	responseMode := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", llm["response_mode"])))
	if responseMode == "streaming" || responseMode == "sse" {
		return true
	}
	if stream, ok := llm["stream"].(bool); ok && stream {
		return true
	}
	return strings.Contains(strings.ToLower(c.GetHeader("Accept")), "text/event-stream")
}

func (autoApi *AutoCodeApi) proxyLLMStream(c *gin.Context, llm common.JSONMap) error {
	res, err := autoCodeService.LLMAutoStream(c.Request.Context(), llm)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("upstream LLM streaming service returned non-2xx: status=%d content-type=%s read-body-err=%w", res.StatusCode, res.Header.Get("Content-Type"), readErr)
		}
		return fmt.Errorf("upstream LLM streaming service returned non-2xx: status=%d content-type=%s body=%s", res.StatusCode, res.Header.Get("Content-Type"), previewResponseBody(body))
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return errors.New("current response does not support streaming output")
	}

	copyLLMStreamHeaders(c.Writer.Header(), res.Header)
	if c.Writer.Header().Get("Content-Type") == "" {
		c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	}
	if c.Writer.Header().Get("Cache-Control") == "" {
		c.Writer.Header().Set("Cache-Control", "no-cache")
	}
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(res.StatusCode)
	flusher.Flush()

	buf := make([]byte, 32*1024)
	for {
		n, readErr := res.Body.Read(buf)
		if n > 0 {
			if _, writeErr := c.Writer.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write streaming response to client: %w", writeErr)
			}
			flusher.Flush()
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return fmt.Errorf("failed to read upstream streaming response: %w", readErr)
		}
	}
}

func copyLLMStreamHeaders(dst, src http.Header) {
	for _, key := range []string{
		"Content-Type",
		"Cache-Control",
		"Content-Encoding",
		"Content-Language",
		"X-Accel-Buffering",
	} {
		if value := src.Get(key); value != "" {
			dst.Set(key, value)
		}
	}
}

func writeLLMStreamError(c *gin.Context, err error) {
	payload, marshalErr := json.Marshal(gin.H{
		"message": err.Error(),
	})
	if marshalErr != nil {
		payload = []byte(`{"message":"streaming proxy failed"}`)
	}
	_, _ = c.Writer.WriteString("event: error\n")
	_, _ = c.Writer.WriteString("data: ")
	_, _ = c.Writer.Write(payload)
	_, _ = c.Writer.WriteString("\n\n")
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
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
