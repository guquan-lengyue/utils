package proxypool

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
)

// 稻壳代理 https://www.docip.net/free
const DocipUrl = "https://www2.docip.net/data/free.json"

type docipProxyResp struct {
	Time uint `json:"time"`
	Data []struct {
		Ip        string `json:"ip"`
		Addr      string `json:"addr"`
		ProxyType string `json:"proxy_type"`
	}
}

func DocipProxy() []*Proxy {
	request, err := httpRequest(DocipUrl, "GET")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	cli := http.Client{Timeout: 60 * time.Second}
	resp, err := cli.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	r := docipProxyResp{}
	err = sonic.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
		return nil
	}
	updateTime := time.Unix(int64(r.Time), 0)
	rst := make([]*Proxy, 0, len(r.Data))
	for _, p := range r.Data {
		var pt ProxyType = HTTP
		if p.ProxyType == "1" {
			pt = HTTPS
		}
		rp := &Proxy{
			Host:   p.Ip,
			Type:   pt,
			Update: updateTime,
			Addr:   p.Addr,
		}
		rst = append(rst, rp)
	}
	return rst
}

func httpRequest(url string, method string, headers ...map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for _, header := range headers {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}
