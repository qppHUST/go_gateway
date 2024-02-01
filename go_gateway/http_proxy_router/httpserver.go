package http_proxy_router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qppHUST/go_gateway/cert_file"
	"github.com/qppHUST/go_gateway/golang_common/lib"
	"github.com/qppHUST/go_gateway/middleware"
)

var (
	HttpServer  *http.Server
	HttpsServer *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpServer = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] http_proxy_run %s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] http_proxy_run %s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

// 仅仅提供了https的代理启动方式
func HttpsServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpsServer = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] https_proxy_run %s\n", lib.GetStringConf("proxy.https.addr"))
	//todo 以下命令只在编译机有效，如果是交叉编译情况下需要单独设置路径
	//if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err!=http.ErrServerClosed {
	if err := HttpsServer.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] https_proxy_run %s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpServer.Shutdown(ctx); err != nil {
		log.Printf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop %v stopped\n", lib.GetStringConf("proxy.http.addr"))
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsServer.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop %v stopped\n", lib.GetStringConf("proxy.https.addr"))
}
