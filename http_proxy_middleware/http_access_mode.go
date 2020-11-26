package http_proxy_middleware

import (
	"fmt"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/gin-gonic/gin"
)

//匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			//未匹配到服务
			middleware.ResponseError(c, middleware.ServiceNotMatchCode, err)
			c.Abort()
			return
		}
		fmt.Println("matched service", public.Obj2Json(service))
		//把服务信息 存到gin的上下文  以方便后面的中间件取服务信息
		c.Set("service", service)
		c.Next()
	}
}
