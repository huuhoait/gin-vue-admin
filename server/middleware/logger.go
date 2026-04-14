package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLayout log layout
type LogLayout struct {
	Time      time.Time
	Metadata  map[string]interface{} // stores custom metadata
	Path      string                 // access path
	Query     string                 // query parameters
	Body      string                 // request body data
	IP        string                 // IP address
	UserAgent string                 // user agent
	Error     string                 // error
	Cost      time.Duration          // time cost
	Source    string                 // source
}

type Logger struct {
	// Filter user-defined filter
	Filter func(c *gin.Context) bool
	// FilterKeyword keyword filter (key)
	FilterKeyword func(layout *LogLayout) bool
	// AuthProcess authentication processing
	AuthProcess func(c *gin.Context, layout *LogLayout)
	// Print log processing
	Print func(LogLayout)
	// Source service unique identifier
	Source string
}

func (l Logger) SetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		var body []byte
		if l.Filter != nil && !l.Filter(c) {
			body, _ = c.GetRawData()
			// put the original body back
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		c.Next()
		cost := time.Since(start)
		layout := LogLayout{
			Time:      time.Now(),
			Path:      path,
			Query:     query,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Error:     strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n"),
			Cost:      cost,
			Source:    l.Source,
		}
		if l.Filter != nil && !l.Filter(c) {
			layout.Body = string(body)
		}
		if l.AuthProcess != nil {
			// process authentication information
			l.AuthProcess(c, &layout)
		}
		if l.FilterKeyword != nil {
			// handle key/value masking etc.
			l.FilterKeyword(&layout)
		}
		// handle log output
		l.Print(layout)
	}
}

func DefaultLogger() gin.HandlerFunc {
	return Logger{
		Print: func(layout LogLayout) {
			// standard output, collected by k8s
			v, _ := json.Marshal(layout)
			fmt.Println(string(v))
		},
		Source: "GVA",
	}.SetLoggerMiddleware()
}
