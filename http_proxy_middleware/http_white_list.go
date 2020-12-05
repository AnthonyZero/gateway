package http_proxy_middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/gin-gonic/gin"
)

//IP白名单
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, middleware.ServiceNotMatchCode, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		iplist := []string{} //白名单
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		//如果开启权限验证 并且 白名单列表 不为空
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			//如果客户端ip 不在白名单中
			if !public.InStringSlice(iplist, c.ClientIP()) {
				middleware.ResponseError(c, middleware.GATEWAY_ERROR_CODE, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
