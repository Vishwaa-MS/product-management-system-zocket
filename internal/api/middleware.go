package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// LoggingMiddleware logs details about each incoming request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		
		logrus.WithFields(logrus.Fields{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"status":       c.Writer.Status(),
			"client_ip":    c.ClientIP(),
			"latency":      duration.String(),
			"user_agent":   c.Request.UserAgent(),
		}).Info("Incoming Request")
	}
}

// AuthMiddleware handles basic authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		username, password, hasAuth := c.Request.BasicAuth()
		
		if !hasAuth || !validateCredentials(username, password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		
		// Set user context for further use
		c.Set("user_id", getUserID(username))
		c.Next()
	}
}

// RateLimitMiddleware prevents too many requests from a single IP
func RateLimitMiddleware() gin.HandlerFunc {
	// Create a rate limiter that allows 100 requests per minute
	limiter := rate.NewLimiter(rate.Limit(100), 100)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check if the request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"details": fmt.Sprintf("Too many requests from %s", clientIP),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ErrorHandlingMiddleware handles panics and unexpected errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				logrus.WithFields(logrus.Fields{
					"error":   err,
					"method":  c.Request.Method,
					"path":    c.Request.URL.Path,
					"client":  c.ClientIP(),
				}).Error("Panic recovered")

				// Respond with internal server error
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"message": "An unexpected error occurred",
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS Middleware to handle Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper functions (these would typically be in a separate authentication service)
func validateCredentials(username, password string) bool {

	return username == "admin" && password == "password"
}

func getUserID(username string) uint {
	// Example: Return a dummy user ID
	if username == "admin" {
		return 1
	}
	return 0
}