# gateway
API gateway(Golang)

## 启动Dashboard(后台管理系统API)
go run main.go -endpoint=dashboard -config=./conf/dev/

## 启动代理服务器
go run main.go -endpoint=server -config=./conf/dev/
