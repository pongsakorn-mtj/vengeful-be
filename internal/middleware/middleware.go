package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  60,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiterCtx, err := instance.Get(c, ip)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Rate limit error",
			})
			c.Abort()
			return
		}

		if limiterCtx.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("Request processed")
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}
