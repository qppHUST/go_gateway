package http_proxy_middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qppHUST/go_gateway/dao"
	"github.com/qppHUST/go_gateway/middleware"
	"github.com/qppHUST/go_gateway/public"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			// 不是租户，直接下一个
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		//用户打开了日请求量上限，同时此时的请求数量已经大于了这个限制
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount)))
			c.Abort()
			return
		}
		c.Next()
	}
}
