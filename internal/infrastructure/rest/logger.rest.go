package rest

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func RestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// Create a custom response writer to capture the response
		buf := new(bytes.Buffer)
		writer := &customResponseWriter{body: buf, ResponseWriter: c.Writer}

		// Replace the original response writer with the custom writer
		c.Writer = writer

		// Proceed with the next middleware or handler
		c.Next()

		// Get the captured response body from the buffer
		responseBody := buf.String()

		if len(c.Errors) > 0 {
			logging.Error(logging.LoggerInput{
				Data: map[string]any{
					"Method":       c.Request.Method,
					"Path":         c.Request.URL.Path,
					"StatusCode":   c.Writer.Status(),
					"ResponseBody": responseBody,
					"ClientIP":     c.ClientIP(),
					"ErrorMessage": c.Errors.String(),
					"Body":         c.Request.Body,
					"Header":       c.Request.Header,
					"Latency":      time.Since(startTime).String(),
				},
			})
		} else {
			logging.Info(logging.LoggerInput{
				Data: map[string]any{
					"Method":       c.Request.Method,
					"Path":         c.Request.URL.Path,
					"StatusCode":   c.Writer.Status(),
					"ResponseBody": responseBody,
					"ClientIP":     c.ClientIP(),
					"ErrorMessage": c.Errors.String(),
					"Body":         c.Request.Body,
					"Header":       c.Request.Header,
					"Latency":      time.Since(startTime).String(),
				},
			})
		}
	}
}

// customResponseWriter is a custom wrapper for gin.ResponseWriter
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write overrides the Write method to capture the response body
func (w *customResponseWriter) Write(b []byte) (int, error) {
	// Write the response to the buffer
	w.body.Write(b)
	// Write the response to the original ResponseWriter
	return w.ResponseWriter.Write(b)
}
