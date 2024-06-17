package proxypool

import (
	"cmp"
	"crypto/tls"
	"math"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"
)

type Register func() []*Proxy

var registers []Register

func registProxyPool(r Register) {
	registers = append(registers, r)
}

// pool
// {ip:port:Proxy}
type pool struct {
	pool []*Proxy
}

var poolInst *pool
var once sync.Once

func GetProxyPool() ProxyPool {
	if poolInst == nil {
		once.Do(func() {
			poolInst = &pool{}
			for _, r := range registers {
				proxy := r()
				poolInst.pool = append(poolInst.pool, proxy...)
			}
		})
	}
	defer task()
	return poolInst
}

func task() {
	go func() {
		for {
			time.Sleep(12 * time.Hour)
			for _, r := range registers {
				proxy := r()
				poolInst.pool = append(poolInst.pool, proxy...)
			}
		}
	}()
	go func() {
		for {
			poolInst.FlashDelay()
			time.Sleep(5 * time.Minute)
		}
	}()
}

var idx = 0
var lock sync.Mutex

func (p *pool) Proxy() (*url.URL, error) {
	lock.Lock()
	defer lock.Unlock()
	proxy := p.pool[idx]
	idx++
	idx %= 10
	return proxy.getUrl()
}

func (p *pool) FlashDelay(targetHostOptions ...string) {
	th := "https://www.baidu.com"
	for _, t := range targetHostOptions {
		th = t
	}
	var wg sync.WaitGroup
	for _, p := range p.pool {
		wg.Add(1)
		go func(proxy *Proxy) {
			defer wg.Done()
			pu, err := proxy.getUrl()
			if err != nil {
				proxy.Delay = math.MaxInt64
				return
			}
			cli := http.Client{
				Transport: &http.Transport{
					Proxy:           http.ProxyURL(pu),
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: 3 * time.Second,
			}
			begin := time.Now()
			resp, err := cli.Get(th)
			end := time.Now()
			if err != nil {
				proxy.Delay = math.MaxInt64
				return
			}
			proxy.Delay = end.Sub(begin).Milliseconds()
			defer resp.Body.Close()

		}(p)
	}
	wg.Wait()
	slices.SortFunc(poolInst.pool, func(a *Proxy, b *Proxy) int {
		return cmp.Compare(a.Delay, b.Delay)
	})
}
