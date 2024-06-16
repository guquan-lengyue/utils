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

const getIp = "https://ipcalf.com"

func Test_Proxy(t *testing.T) {
	proxyPool := DocipProxy()
	wg := sync.WaitGroup{}
	pass := true
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
			if err != nil {
				return
			}
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
			if err != nil {
				return
			}
			proxy.Delay = end.Sub(begin).Milliseconds()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			hostIdx := strings.Index(proxy.Host, ":")
			ipIsContain := bytes.Contains(body, []byte(proxy.Host[:hostIdx]))
			t.Logf("成功 %#v, %v", proxy, ipIsContain)
			if !pass {
				pass = true
			}
		}(proxy)
	}
	wg.Wait()
	assert.True(t, pass)
}
