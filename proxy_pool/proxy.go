package proxypool

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const supportorUrl = "http://api.89ip.cn/tqdl.html?api=1&num=60&port=&address=&isp="

var proxyReg = regexp.MustCompile(`(.*?):(\d+)<br>`)

type proxyPool struct {
	host string
	port int
}

func getProxyHost() ([]*proxyPool, error) {
	req, err := http.NewRequest("GET", supportorUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	cli := &http.Client{Timeout: 60 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	proxy := proxyReg.FindAllSubmatch(body, -1)
	rst := make([]*proxyPool, 0, len(proxy))
	for i := 0; i < len(proxy); i++ {
		p, err := strconv.Atoi(string(proxy[i][2]))
		if err != nil {
			continue
		}
		rst = append(rst, &proxyPool{host: string(proxy[i][1]), port: p})
	}
	return rst, nil
}
