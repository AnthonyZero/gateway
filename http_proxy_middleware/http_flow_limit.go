package http_proxy_middleware

import (
	"errors"
	"fmt"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/gin-gonic/gin"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, middleware.ServiceNotMatchCode, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		//根据access_control 表的字段来限流
		//服务端限流
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				middleware.ResponseError(c, middleware.InternalErrorCode, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				middleware.ResponseError(c, middleware.REQUEST_LIMITED_CODE, errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		//客户端限流 基于clientIP限流
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				middleware.ResponseError(c, middleware.InternalErrorCode, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, middleware.REQUEST_LIMITED_CODE, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
