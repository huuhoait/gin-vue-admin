package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TimeoutMiddleware creates a timeout middleware
// parameter timeout sets the timeout duration (e.g., time.Second * 5)
// usage example: xxx.Get("path",middleware.TimeoutMiddleware(30*time.Second),HandleFunc)
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// use buffered channel to avoid goroutine leak
		done := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					select {
					case panicChan <- p:
					default:
					}
				}
				select {
				case done <- struct{}{}:
				default:
				}
			}()
			c.Next()
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			return
		case <-ctx.Done():
			// ensure server timeout is set long enough
			c.Header("Connection", "close")
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"code": 504,
				"msg":  "request timeout",
			})
			return
		}
	}
}
