package proxypool

import "time"

type ProxyType uint

const (
	HTTP ProxyType = iota
	HTTPS
)

type Proxy struct {
	Host   string        // 代理地址
	Addr   string        // ip地址位置
	Type   ProxyType     // 协议
	Update time.Time     // 更新时间
	Delay  time.Duration // 目标延迟
}
