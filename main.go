package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthonyzero/gateway/http_proxy_router"

	"github.com/anthonyzero/gateway/router"
	"github.com/e421083458/golang_common/lib"
)

//endpoint:  dashboard ->后台管理接口  server->代理服务器
//config:  ./conf/prod/ 对应配置文件夹

//终端输入参数
var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		//自动打印错误提示
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		//自动打印错误提示
		flag.Usage()
		os.Exit(1)
	}

	//要么运行后台接口部分 要么运行代理服务器
	if *endpoint == "dashboard" {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()

		//启动各种代理服务
		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}
}
