package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/gopkg/server/xgin"
	"net/url"
)

func ProxyExample(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var leader string
		if leader != "" {
			f, _ := ProxyServer()
			//utils.SetCtxLog(ctx, fmt.Sprintf("need forward=%s", addr))
			url, err := url.ParseRequestURI(fmt.Sprintf("http://%s%s", leader, ctx.Request.RequestURI))
			if err != nil {
				xgin.ResponseNo(ctx, xgin.ErrorForwarderUrl, fmt.Sprintf("forwarder parse url=%s failed cause=%s", url, leader))
				return
			}
			ctx.Request.URL = url
			f.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}
		next(ctx)
	}
}
