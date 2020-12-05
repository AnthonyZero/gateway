package http_proxy_middleware

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/gin-gonic/gin"
)

//url重写  ^/test_http_service/abb/(.*) /test_http_service/bba/$1
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, middleware.ServiceNotMatchCode, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			//fmt.Println("item rewrite",item)
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err", err)
				continue
			}
			//fmt.Println("before rewrite",c.Request.URL.Path)
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			//fmt.Println("after rewrite",c.Request.URL.Path)
		}
		c.Next()
	}
}
