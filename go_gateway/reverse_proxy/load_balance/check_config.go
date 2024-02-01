package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	//default check setting
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

// 主动探测方式
type LoadBalanceCheckConf struct {
	observers    []Observer        //负载均衡器
	confIpWeight map[string]string //ip 权重映射
	activeList   []string          //活跃的ip列表
	format       string            //生成下游地址+端口的格式化字符串
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// 获取配置的ip列表
func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50" //默认weight
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// 心跳检测
func (s *LoadBalanceCheckConf) WatchConf() {
	//fmt.Println("watchConf")
	go func() {
		confIpErrNum := map[string]int{}
		for {
			// 以一个for循环来实现对可用ip的查找
			changedList := []string{}
			for item := range s.confIpWeight {
				//测试哪些ip是可用的
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				//todo http statuscode
				if err == nil {
					conn.Close()
					//此ip可用，出错次数归零
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				if err != nil {
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] += 1
					} else {
						confIpErrNum[item] = 1
					}
				}
				//出错次数少于规定的次数，就把这个ip放进changedlist
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			//当前活跃ip列表和改变之后的ip列表不一样
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

// UpdateConf 更新配置
func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	//fmt.Println("UpdateConf", conf)
	s.activeList = conf
	for _, obs := range s.observers {
		// 更新配置时，通知监听者也更新
		obs.Update()
	}
}

// NewLoadBalanceCheckConf 返回配置负载均衡器配置
func NewLoadBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConf, error) {
	aList := []string{}
	//默认初始化
	for item := range conf {
		aList = append(aList, item)
	}
	mConf := &LoadBalanceCheckConf{format: format, activeList: aList, confIpWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}
