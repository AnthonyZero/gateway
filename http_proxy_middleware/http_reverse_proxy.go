package http_proxy_middleware

import (
	"errors"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/reverse_proxy"

	"github.com/anthonyzero/gateway/middleware"
	"github.com/gin-gonic/gin"
)

//HTTP反向代理中间件
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从gin的上下文获取服务详情
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, middleware.GATEWAY_ERROR_CODE, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail)
		//获取负载均衡器
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, middleware.GATEWAY_ERROR_CODE, err)
			c.Abort()
			return
		}
		//每个HTTP服务 使用单独的transport连接池
		trans, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, middleware.GATEWAY_ERROR_CODE, err)
			c.Abort()
			return
		}
		//middleware.ResponseSuccess(c,"ok")
		//return
		//创建 reverseproxy
		//使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
