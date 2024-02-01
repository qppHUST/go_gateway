package middleware

import (
	"fmt"

	"github.com/qppHUST/go_gateway/golang_common/lib"
	"github.com/gin-gonic/gin"
)

// 管理后台的IP名单
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range lib.GetStringSliceConf("base.http.allow_ip") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(c, InternalErrorCode, fmt.Errorf("%v, not in iplist", c.ClientIP()))
			c.Abort()
			return
		}
		c.Next()
	}
}
