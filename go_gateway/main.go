package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/qppHUST/go_gateway/dao"
	"github.com/qppHUST/go_gateway/golang_common/lib"
	"github.com/qppHUST/go_gateway/grpc_proxy_router"
	"github.com/qppHUST/go_gateway/http_proxy_router"
	"github.com/qppHUST/go_gateway/router"
	"github.com/qppHUST/go_gateway/tcp_proxy_router"
)

//endpoint dashboard后台管理  server代理服务器
//config ./conf/prod/ 对应配置文件夹

// 使用flag包读取命令行变量
var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		//打印endpoint的用法
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		//后台管理
		lib.InitModule(*config)
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		//网关代理
		lib.InitModule(*config)
		defer lib.Destroy()
		//进行一些初始化
		dao.ServiceManagerHandler.LoadOnce()
		dao.AppManagerHandler.LoadOnce()

		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()
		go func() {
			tcp_proxy_router.TcpServerRun()
		}()
		go func() {
			grpc_proxy_router.GrpcServerRun()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpServerStop()
		grpc_proxy_router.GrpcServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}
}
