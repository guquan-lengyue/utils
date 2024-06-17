package proxypool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_InitProxyPool(t *testing.T) {
	a := GetProxyPool()
	assert.NotNil(t, a)
	time.Sleep(4 * time.Second)
	u, err := a.Proxy()
	assert.NoError(t, err)
	t.Log(u.Host)
}
