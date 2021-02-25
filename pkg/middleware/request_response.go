package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

const ReservedBodyMap = "REQUEST_BODY_MAP"
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

//now only for json body 和 x-www-form-urlencoded body
func ParseRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != http.NoBody {
			var buf bytes.Buffer
			contentType := c.Request.Header.Get("content-type")
			if _, err := buf.ReadFrom(c.Request.Body); err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
				switch strings.ToLower(contentType) {
				//todo const
				case "application/json":
					params := make(map[string]interface{})
					if err = json.Unmarshal(buf.Bytes(), &params); err == nil {
						c.Set(ReservedBodyMap, params)
					}
				case "application/x-www-form-urlencoded":
					params := make(map[string]interface{})
					for _, param := range strings.Split(buf.String(), "&") {
						value := strings.Split(param, "=")
						if len(value) == 2 {
							params[value[0]] = value[1]
						}
					}
					c.Set(ReservedBodyMap, params)
				default:
					//todo log 不支持

				}
			}

		}
		c.Next()
	}
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
