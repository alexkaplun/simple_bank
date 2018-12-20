package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func ZapLogger(logger *zap.Logger) gin.HandlerFunc{
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// after request
		logger.Info(c.Request.URL.Path,
			zap.String("method", c.Request.Method),
			zap.String("remote_addr", c.Request.RemoteAddr),
			zap.String("url", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("work_time", time.Since(start)),
		)
	}
}