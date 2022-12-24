package xgin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/gopkg/log"
	"go.uber.org/zap"
	"runtime"
	"time"
)

const (
	RequestTime   = "request_Time"
	LoggerInfo    = "logger_info"
	ErrorResponse = "error_response"
)

func SetErrorResponse(c *gin.Context) {
	c.Set(ErrorResponse, ErrorResponse)
}

func LoggerFormatTools(c *gin.Context, info string) {
	v, ok := c.Get(LoggerInfo)
	if ok {
		str, ok := v.(string)
		if ok {
			str = str + info
			c.Set(LoggerInfo, str)
		}
	} else {
		c.Set(LoggerInfo, info)
	}
}

func LoggerJsonTools(c *gin.Context, field ...zap.Field) {
	if len(field) == 0 {
		return
	}
	v, ok := c.Get(LoggerInfo)
	if ok {
		zf, ok := v.([]zap.Field)
		if ok {
			zf = append(zf, field...)
			c.Set(LoggerInfo, zf)
		}
	} else {
		c.Set(LoggerInfo, field)
	}
}

func RecordFormat(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestTime, time.Now())
		defer func() {
			var panic string
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				panic = fmt.Sprintf("panic Error: %v;stack: %s", err, buf[:n])
				ResponseNo(c, ErrorUnknown)
			}
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("received from %s %s rt=%d ", c.RemoteIP(), c.Request.RequestURI, time.Since(c.MustGet(RequestTime).(time.Time)).Milliseconds()))
			v, ok := c.Get(LoggerInfo)
			if ok {
				str, ok := v.(string)
				if ok && str != "" {
					buf.WriteString(str)
				}
			}
			_, e := c.Get(ErrorResponse)
			if panic != "" || e {
				buf.WriteString(fmt.Sprintf("panic=%s ", panic))
				log.Logger().Error(buf.String())
				return
			}
			log.Logger().Info(buf.String())
		}()
		next(c)
	}
}

func RecordJson(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestTime, time.Now())
		defer func() {
			var panic string
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				panic = fmt.Sprintf("panic Error: %v;stack: %s", err, buf[:n])
				ResponseNo(c, ErrorUnknown)
			}

			var zf []zap.Field
			zf = append(zf,
				zap.String("ip", c.RemoteIP()),
				zap.String("uri", c.Request.RequestURI),
				zap.Int64("rt", time.Since(c.MustGet(RequestTime).(time.Time)).Milliseconds()),
			)
			v, ok := c.Get(LoggerInfo)
			if ok {
				fs, ok := v.([]zap.Field)
				if ok && len(fs) != 0 {
					zf = append(zf, fs...)
				}
			}
			_, e := c.Get(ErrorResponse)
			if panic != "" || e {
				zf = append(zf, zap.String("panic", panic))
				log.LoggerJ().Error("request received", zf...)
				return
			}
			log.LoggerJ().Info("request received", zf...)
		}()
		next(c)
	}
}
