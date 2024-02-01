package http_proxy_router

import (
	"github.com/qppHUST/go_gateway/controller"
	"github.com/qppHUST/go_gateway/http_proxy_middleware"
	"github.com/qppHUST/go_gateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	//todo 优化点1 ，default会打印debug信息
	//router := gin.Default()
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//合并部署方式，将前端整合到后端来
	// router.Static("/dist", "./dist")

	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}

	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),     //匹配服务（根据域名或者前缀，load对应的serviceDetail）
		http_proxy_middleware.HTTPFlowCountMiddleware(),      //请求数量统计
		http_proxy_middleware.HTTPFlowLimitMiddleware(),      //请求限流
		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),   //jwt认证
		http_proxy_middleware.HTTPJwtFlowCountMiddleware(),   //jwt请求统计
		http_proxy_middleware.HTTPJwtFlowLimitMiddleware(),   //jwt请求限流
		http_proxy_middleware.HTTPWhiteListMiddleware(),      //白名单
		http_proxy_middleware.HTTPBlackListMiddleware(),      //黑名单
		http_proxy_middleware.HTTPHeaderTransferMiddleware(), //http header 修改
		http_proxy_middleware.HTTPStripUriMiddleware(),       //跳过uri
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),     //url重写
		http_proxy_middleware.HTTPReverseProxyMiddleware())   //代理

	return router
}
