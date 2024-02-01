package http_proxy_middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qppHUST/go_gateway/dao"
	"github.com/qppHUST/go_gateway/middleware"
	"github.com/qppHUST/go_gateway/public"
)

// 根据请求类型，保存对应的serviceDetail
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		fmt.Println("matched service", public.Obj2Json(service))
		c.Set("service", service)
		c.Next()
	}
}
