package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type ResponseCode int

//1000以下为通用码，1000以上为用户自定义码
const (
	SuccessCode ResponseCode = iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	BusinessErrorCode       ResponseCode = 400
	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000
	ParamCheckErrorCode     ResponseCode = 1001
	RecordExistCode         ResponseCode = 1002
	ServiceNotMatchCode     ResponseCode = 1003
	GROUPALL_SAVE_FLOWERROR ResponseCode = 2001
	REVERSE_PROXY_ERROR     ResponseCode = 9999
)

type Response struct {
	Code     ResponseCode `json:"code"`
	ErrorMsg string       `json:"errmsg"`
	Data     interface{}  `json:"data"`
	TraceId  interface{}  `json:"trace_id"`
	Stack    interface{}  `json:"stack"`
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	stack := ""
	if c.Query("is_debug") == "1" || lib.GetConfEnv() == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}

	resp := &Response{Code: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId, Stack: stack}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	resp := &Response{Code: SuccessCode, ErrorMsg: "", Data: data, TraceId: traceId}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
