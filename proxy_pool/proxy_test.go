package proxypool

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetProxy(t *testing.T) {
	pool := DocipProxy()
	assert.NotEmpty(t, pool, "代理池不应该为空")
}

const getIp = "http://www.baidu.com"

func Test_Proxy(t *testing.T) {
	proxyPool := DocipProxy()
	wg := sync.WaitGroup{}
	for _, proxy := range proxyPool {
		wg.Add(1)
		go func(proxy *Proxy) {
			defer func() {
				wg.Done()
				if e := recover(); e != nil {
					// ignore
				}
			}()
			protocol := "http"
			if proxy.Type == HTTPS {
				protocol = "https"
			}
			proxyUrl, err := url.Parse(protocol + "://" + proxy.Host)
			assert.NoError(t, err)
			cli := http.Client{
				Transport: &http.Transport{
					Proxy:           http.ProxyURL(proxyUrl),
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: 10 * time.Second,
			}
			begin := time.Now()
			resp, err := cli.Get(getIp)
			end := time.Now()
			assert.NoError(t, err, "请求错误,使用代理："+proxy.Host+"地址："+proxy.Addr)
			proxy.Delay = end.Sub(begin).Milliseconds()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			hostIdx := strings.Index(proxy.Host, ":")
			ipIsContain := bytes.Contains(body, []byte(proxy.Host[:hostIdx]))
			// assert.True(t, ipIsContain)

			t.Logf("成功 %#v, %v", proxy, ipIsContain)
		}(proxy)
	}
	wg.Wait()
}
 