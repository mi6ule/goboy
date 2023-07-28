package rest

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func RequestLoggerMiddleware(param gin.LogFormatterParams) string {
	if param.ErrorMessage == "" {
		logging.Info(logging.LoggerInput{
			Data: map[string]any{
				"ClientIP":   param.ClientIP,
				"Method":     param.Method,
				"Path":       param.Path,
				"StatusCode": param.StatusCode,
				"Latency":    param.Latency,
				"Body":       param.Request.Body,
				"Header":     param.Request.Header,
			},
		})
	} else {
		logging.Error(logging.LoggerInput{
			Data: map[string]any{
				"ClientIP":     param.ClientIP,
				"Method":       param.Method,
				"Path":         param.Path,
				"StatusCode":   param.StatusCode,
				"Latency":      param.Latency,
				"ErrorMessage": param.ErrorMessage,
				"Body":         param.Request.Body,
				"Header":       param.Request.Header,
			},
		})
	}
	return ""
}

// LogResponseMiddleware is a custom middleware to log the response body
func LogResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a custom response writer to capture the response
		buf := new(bytes.Buffer)
		writer := &customResponseWriter{body: buf, ResponseWriter: c.Writer}

		// Replace the original response writer with the custom writer
		c.Writer = writer

		// Proceed with the next middleware or handler
		c.Next()

		// Get the captured response body from the buffer
		responseBody := buf.String()

		// Log the response body
		logging.Info(logging.LoggerInput{
			Data: map[string]any{
				"Method":       c.Request.Method,
				"Path":         c.Request.URL.Path,
				"StatusCode":   c.Writer.Status(),
				"ResponseBody": responseBody,
			},
		})
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
