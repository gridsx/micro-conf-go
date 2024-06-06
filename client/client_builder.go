package client

import (
	"sync/atomic"
)

func NewClient(serviceName, group, metaServers, token string, filters, headers map[string]string,
	refreshPeriod, timeout, rateLimit int) *Client {
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
	err := svc.init(metaServers, token)
	if err != nil {
		return nil
	}
	return svc
}
