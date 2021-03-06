package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	. "server/pkg/config"
	"server/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

func init() {
	if err := initLogger("info", GlobalConfig.Accesslog, 500, 10, 10); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
}

func initLogger(level string, filename string, maxSize, maxBackup, maxAge int) (err error) {
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	Logger = zap.New(core, zap.AddCaller())
	SugarLogger = Logger.Sugar()
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件路径
		MaxSize:    maxSize,   // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackup, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,    // 文件最多保存多少天
		Compress:   true,      // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		requestBody, responseBody := "", ""
		if v, ok := c.Get(middleware.RequestBody); ok {
			requestBody = v.(string)
		}
		if v, ok := c.Get(middleware.ResponseBody); ok {
			responseBody = v.(string)
		}
		cost := time.Since(start) * 1000
		Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("requestBody", requestBody),
			zap.String("responseBody", responseBody),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志，并推送钉钉消息
func GinRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 打印错误日志
		httpRequest, _ := httputil.DumpRequest(c.Request, true)
		request := string(httpRequest)
		stack := string(debug.Stack())
		Logger.Error("[Recovery from panic]",
			zap.Any("error", recovered),
			zap.String("request", request),
			zap.String("stack", stack),
		)
		if GlobalConfig.Env == "local" || GlobalConfig.Env == "dev" || GlobalConfig.Env == "test" {
			fmt.Println("\n[GinRecovery->error]", recovered, stack)
		}

		// 推送钉钉
		if GlobalConfig.Env == "prod" && GlobalConfig.ErrorDingTalk != "" {
			var s = fmt.Sprintf("[%s][ERROR] Request: %s Stack: %s", GlobalConfig.Env, request, stack)
			DingTalk(s)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func DingTalk(s string) {
	content, data := make(map[string]string), make(map[string]interface{})
	content["content"] = s

	data["msgtype"] = "text"
	data["text"] = content
	b, _ := json.Marshal(data)

	resp, err := http.Post(GlobalConfig.ErrorDingTalk,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func Debug(args ...interface{}) {
	SugarLogger.Debug(args...)
}

func Info(args ...interface{}) {
	SugarLogger.Info(args...)
}

func Warn(args ...interface{}) {
	SugarLogger.Warn(args...)
}

func Error(args ...interface{}) {
	SugarLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	SugarLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	SugarLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	SugarLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	SugarLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	SugarLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	SugarLogger.Fatalf(template, args...)
}
