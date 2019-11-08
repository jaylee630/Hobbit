package ctlmain

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type respDataWriter struct {
	gin.ResponseWriter
	data *bytes.Buffer
}

func (w respDataWriter) Write(b []byte) (int, error) {
	w.data.Write(b)
	return w.ResponseWriter.Write(b)
}

// RespLoggerWithWriter middleware log response data when status code is larger(include equal) than 400.
// But this is not a prettiest way of doing this, because it cache all resp no matter status code is.
// You SHOULD NOT use it at production environment.
func RespLoggerWithWriter(out io.Writer) gin.HandlerFunc {

	return func(c *gin.Context) {

		path := c.Request.URL.Path

		w := respDataWriter{data: bytes.NewBuffer(make([]byte, 0)), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		if c.Writer.Status() < 400 {
			return
		}

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		end := time.Now()

		_, _ = fmt.Fprintf(out, "[GIN] %v | %3d | %13s | %-7s  %s [response]\n%s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			clientIP,
			method,
			path,
			w.data.String(),
		)
	}
}

func NewGinEngine(out io.Writer, logLevel string) *gin.Engine {

	// init gin engine
	engine := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// init gin debug functions
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		engine.Use(RespLoggerWithWriter(out))
	}

	// init log and recovery
	engine.Use(gin.RecoveryWithWriter(out))
	engine.Use(gin.LoggerWithWriter(out))

	return engine
}
