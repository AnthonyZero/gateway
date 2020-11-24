package dao

import (
	"net/http/httptest"
	"sync"

	"github.com/anthonyzero/gateway/dto"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

//服务管理
var ServiceManagerHandler *ServiceManager

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail //服务名称 -》服务详情 通过服务名称的方式取
	ServiceSlice []*ServiceDetail          //通过遍历的方式取
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

//init()函数是在文件包首次被加载的时候执行，且只执行一次
func init() {
	ServiceManagerHandler = NewServiceManager()
}
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

//一次性加载服务信息  sync.Onc 是在代码运行中需要的时候执行，且只执行一次
func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		serviceInfo := &ServiceInfo{}
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.err = err
			return
		}
		//获取所有服务列表
		params := &dto.ServiceListInput{PageNo: 1, PageSize: 99999}
		list, _, err := serviceInfo.PageList(c, tx, params)
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			//获取服务详情
			serviceDetail, err := tmpItem.ServiceDetail(c, tx, &tmpItem)
			//fmt.Println("serviceDetail")
			//fmt.Println(public.Obj2Json(serviceDetail))
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.err
}
