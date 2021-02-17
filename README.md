# gateway
API gateway(Golang)

## 启动Dashboard(后台管理系统API)
go run main.go -endpoint=dashboard -config=./conf/dev/
> http://127.0.0.1:8880/swagger/index.html

## 启动代理服务器
go run main.go -endpoint=server -config=./conf/dev/

## 压测
ab -n1000 -c10 http://127.0.0.1:8080/test_http_service/abb

## 通过租户AppId和Secret获取JWT Token
curl 'http://127.0.0.1:8080/oauth/tokens' -u app_id_a:449441eb5e72dca9c42a12f3924ea3a2 -d 'grant_type=client_credentials&scope=read_write'


