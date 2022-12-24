package xgin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	Ok = 0
)

type Page struct {
	Total    int `json:"total"`
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

func (p *Page) JsonString() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Info string      `json:"info"`
}

func NewResponse(code int, info string) Result {
	return Result{Code: code, Info: info}
}

func ResponseOk(c *gin.Context) {
	c.Header("content-type", "application/json")
	c.SecureJSON(http.StatusOK, &Result{Code: Ok, Data: ""})
}

func Response(c *gin.Context, data interface{}) {
	c.Header("content-type", "application/json")
	c.SecureJSON(http.StatusOK, &Result{Code: Ok, Data: data})
}

func ResponseNo(c *gin.Context, data Result, errInfo ...string) {
	c.Header("content-type", "application/json")
	c.SecureJSON(http.StatusOK, &Result{Code: data.Code, Info: data.Info + ":" + strings.Join(errInfo, "")})
}
