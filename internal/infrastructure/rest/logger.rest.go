package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func RestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		buf := new(bytes.Buffer)
		writer := &customResponseWriter{body: buf, ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()
		responseBody := buf.String()
		var reqBodyJSON any
		var resBodyJSON any
		logData := map[string]any{
			"Method":        c.Request.Method,
			"Path":          c.Request.URL.Path,
			"StatusCode":    c.Writer.Status(),
			"ClientIP":      c.ClientIP(),
			"ErrorMessage":  c.Errors.String(),
			"RequestHeader": c.Request.Header,
			"Latency":       time.Since(startTime).String(),
			"ResponseBody":  responseBody,
			"RequestBody":   string(reqBody),
		}
		if err := json.Unmarshal(reqBody, &reqBodyJSON); err == nil {
			logData["RequestBody"] = reqBodyJSON
		}

		if err := json.Unmarshal([]byte(responseBody), &resBodyJSON); err == nil {
			logData["ResponseBody"] = resBodyJSON
		}
		if len(c.Errors) > 0 {
			logData["ErrorMessage"] = c.Errors.String()
			logging.Error(logging.LoggerInput{
				Data: logData,
			})
		} else {
			logging.Info(logging.LoggerInput{
				Data: logData,
			})
		}
	}
}

type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
