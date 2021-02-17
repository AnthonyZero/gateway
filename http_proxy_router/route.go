package http_proxy_router

import (
	"github.com/anthonyzero/gateway/controller"
	"github.com/anthonyzero/gateway/http_proxy_middleware"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	//todo 优化点1
	//router := gin.Default()
	router := gin.New()
	router.Use(middlewares...) //默认使用的中间件
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}

	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),     //匹配到具体服务
		http_proxy_middleware.HTTPFlowCountMiddleware(),      // 计数器
		http_proxy_middleware.HTTPFlowLimitMiddleware(),      // 限流器 令牌桶
		http_proxy_middleware.HTTPWhiteListMiddleware(),      //白名单 （开启了权限验证）
		http_proxy_middleware.HTTPBlackListMiddleware(),      //黑名单 （开启了权限验证）
		http_proxy_middleware.HTTPHeaderTransferMiddleware(), //header transfer add edit等
		http_proxy_middleware.HTTPStripUriMiddleware(),       // 网关处理接入前缀
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),     //url重写
		http_proxy_middleware.HTTPReverseProxyMiddleware())   //代理
	return router
}
