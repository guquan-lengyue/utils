package proxypool

import (
	"testing"
)

func Test_ProxyHost(t *testing.T) {
	rst, err := getProxyHost()
	if err != nil || len(rst) == 0 {
		panic(err)
	}
}
