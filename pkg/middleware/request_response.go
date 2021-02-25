package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const RequestBody = "REQUEST_BODY"
const ResponseBody = "RESPONSE_BODY"

type logWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *logWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// catch requestInfo and responseInfo for log
func CatchRequestAndResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != http.NoBody {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(c.Request.Body); err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
				c.Set(RequestBody, string(buf.Bytes()))
			}

		}
		blw := &logWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		c.Set(ResponseBody, blw.body.String())
	}
}
