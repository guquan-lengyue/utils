package proxypool

import (
	"net/url"
	"time"
)

type ProxyType uint

const (
	HTTP ProxyType = iota
	HTTPS
)

type Proxy struct {
	Host   string    // 代理地址
	Addr   string    // ip地址位置
	Type   ProxyType // 协议
	Update time.Time // 更新时间
	Delay  int64     // 目标延迟
}

func (p *Proxy) getUrl() (*url.URL, error) {
	protocol := "http"
	if p.Type == HTTPS {
		protocol = "https"
	}
	return url.Parse(protocol + "://" + p.Host)
}

type ProxyPool interface {
	Proxy() (*url.URL, error)
	FlashDelay(targetHostOptions ...string)
}
