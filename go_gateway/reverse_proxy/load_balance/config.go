package load_balance

// 配置主题
type LoadBalanceConf interface {
	Attach(o Observer) //加入一个监听者
	GetConf() []string
	WatchConf() //轮询监听下游可用ip的变化
	UpdateConf(conf []string)
}

// Observer 是抽象出来的更新接口，实现它的就是那几个轮询方式
type Observer interface {
	Update()
}
