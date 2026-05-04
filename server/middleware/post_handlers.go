package middleware

import "github.com/gin-gonic/gin"

var postHandlers []gin.HandlerFunc

// RegisterPostHandler appends a hook that runs after the main handler chain
// completes. Intended for plugin init() to observe authenticated requests
// (e.g., online-session tracking) without retrofitting every router group.
// Not thread-safe — call only at process init time.
func RegisterPostHandler(h gin.HandlerFunc) {
	postHandlers = append(postHandlers, h)
}

// PostHandlerChain runs every RegisterPostHandler hook after c.Next(). Hooks
// must not write a response — the response body has already been committed
// to the client by the time they run.
func PostHandlerChain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, h := range postHandlers {
			h(c)
		}
	}
}
