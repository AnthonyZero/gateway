package load_balance

// 配置主题
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

//用于添加ip信息到负载均衡策略
type Observer interface {
	Update()
}
