package http_proxy_middleware

import (
	"errors"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"

	"github.com/gin-gonic/gin"
)

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, middleware.ServiceNotMatchCode, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//统计项 1 全站 2 服务 3 租户
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			middleware.ResponseError(c, middleware.BusinessErrorCode, err)
			c.Abort()
			return
		}
		totalCounter.Increase()

		// dayCount, _ := totalCounter.GetDayData(time.Now())
		// fmt.Printf("totalCounter qps:%v,dayCount:%v", totalCounter.QPS, dayCount)
		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, middleware.BusinessErrorCode, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()

		// dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		// fmt.Printf("serviceCounter qps:%v,dayCount:%v", serviceCounter.QPS, dayServiceCount)
		c.Next()
	}
}
