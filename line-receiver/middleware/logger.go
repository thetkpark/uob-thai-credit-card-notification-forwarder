package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func HttpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		msg := fmt.Sprintf("%s %s %s %d (%d ms)", c.Request.Method, c.Request.URL.Path, c.Request.Proto, c.Writer.Status(), duration.Milliseconds())
		slog.InfoContext(c.Request.Context(), msg)
	}
}
