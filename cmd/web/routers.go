package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())

	// Custom middleware to handle timeouts
	router.Use(timeoutMiddleware(60 * time.Second))

	router.Static("/static", "./static")

	router.POST("api/filter", app.FilterdList)

	router.GET("/", app.ShowHome)
	router.GET("/:page", app.ShowPage)

	return router
}

// Define other functions like timeoutMiddleware and ShowHome as needed

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Create a new channel to signal completion of request handling
		done := make(chan struct{}, 1)

		// Update the context of the request
		c.Request = c.Request.WithContext(ctx)

		// Run the next handler in a goroutine
		go func() {
			c.Next()           // Process request
			done <- struct{}{} // Signal completion
		}()

		// Wait for either completion or timeout
		select {
		case <-ctx.Done():
			// Timeout occurred
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"message": "request timeout"})
		case <-done:
			// Request processed successfully
		}
	}
}
