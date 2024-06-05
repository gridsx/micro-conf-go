package client

import (
	"sync"
	"sync/atomic"
)

var (
	clientMap = sync.Map{}
)

func NewClient(serviceName, group, metaServers string, filters, headers map[string]string,
	refreshPeriod, timeout, rateLimit int) *Client {
	service, ok := clientMap.Load(serviceName)
	if ok && service != nil {
		if svc, ok := service.(*Client); ok {
			return svc
		}
	}
	if group == "" {
		group = "default"
	}
	if refreshPeriod <= 0 {
		refreshPeriod = defaultTimeout / 2
	}
	svc := &Client{
		Name:          serviceName,
		Group:         group,
		Timeout:       timeout,
		RefreshPeriod: refreshPeriod,
		Headers:       headers,
		RateLimit:     rateLimit,
		Filters:       filters,
		instNum:       atomic.Int32{},
	}
	svc.init(metaServers)
	clientMap.Store(serviceName, svc)
	return svc
}

func GetClient(serviceName string) *Client {
	c, ok := clientMap.Load(serviceName)
	if c != nil && ok {
		if client, convertOk := c.(*Client); convertOk {
			return client
		}
	}
	return nil
}
