package dao

import (
	"fmt"
	"strings"
	"sync"

	"github.com/anthonyzero/gateway/public"
	"github.com/anthonyzero/gateway/reverse_proxy/load_balance"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type LoadBalance struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	ServiceID     int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) Find(c *gin.Context, tx *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	model := &LoadBalance{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (t *LoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

//获取负载均衡的IP列表
func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}

//获取权重列表
func (t *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(t.WeightList, ",")
}

//负载均衡器管理
var LoadBalancerHandler *LoadBalancer

type LoadBalancer struct {
	LoadBanlanceMap   map[string]*LoadBalancerItem
	LoadBanlanceSlice []*LoadBalancerItem
	Locker            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBanlance load_balance.LoadBalance
	ServiceName  string
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBanlanceMap:   map[string]*LoadBalancerItem{},
		LoadBanlanceSlice: []*LoadBalancerItem{},
		Locker:            sync.RWMutex{},
	}
}

//init()函数是在文件包首次被加载的时候执行，且只执行一次
func init() {
	LoadBalancerHandler = NewLoadBalancer()
}

//通过服务详情 获取对应的负载均衡策略
func (lbr *LoadBalancer) GetLoadBalancer(service *ServiceDetail) (load_balance.LoadBalance, error) {
	//有直接返回
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			return lbrItem.LoadBanlance, nil
		}
	}
	schema := "http://"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}
	if service.Info.LoadType == public.LoadTypeTCP || service.Info.LoadType == public.LoadTypeGRPC {
		schema = ""
	}
	//获取下游的ip 和权重（如果有）
	ipList := service.LoadBalance.GetIPListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	ipConf := map[string]string{} //ip <-> 权重
	for ipIndex, ipItem := range ipList {
		ipConf[ipItem] = weightList[ipIndex]
	}
	//fmt.Println("ipConf", ipConf)
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}
	//通过类型值和配置 初始化负载均衡策略
	lb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)

	//save to map and slice
	lbItem := &LoadBalancerItem{
		LoadBanlance: lb,
		ServiceName:  service.Info.ServiceName,
	}
	lbr.LoadBanlanceSlice = append(lbr.LoadBanlanceSlice, lbItem)

	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBanlanceMap[service.Info.ServiceName] = lbItem
	return lb, nil
}
