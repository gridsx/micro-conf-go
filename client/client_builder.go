package client

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type clientManager struct {
	baseClientMap   map[string]*Client
	customClientMap map[string]any
	lock            sync.Mutex
}

func (cm *clientManager) NewClient(serviceName, group, metaServers, token string, filters, headers map[string]string,
	refreshPeriod, timeout, rateLimit int) *Client {
	if group == "" {
		group = "default"
	}
	k := fmt.Sprintf("%s_%s", serviceName, group)
	if v, ok := cm.baseClientMap[k]; ok {
		return v
	} else {
		if refreshPeriod <= 0 {
			refreshPeriod = defaultTimeout / 2
		}
		cm.lock.Lock()
		defer cm.lock.Unlock()
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
		cm.baseClientMap[k] = svc
		return svc
	}
}

func (cm *clientManager) MakeClient(serviceName, group string, v any) error {
	if group == "" {
		group = "default"
	}
	k := fmt.Sprintf("%s_%s", serviceName, group)
	if v, ok := cm.baseClientMap[k]; ok {
		if _, ok := cm.customClientMap[k]; ok {
			return errors.New("already made, ensure there is one instance for your client")
		}
		cm.lock.Lock()
		defer cm.lock.Unlock()
		maker := clientMaker{cli: v}
		maker.generateMethods(v)
		cm.customClientMap[k] = v
		return nil
	}
	return errors.New("no base client")
}

var DefaultMgr = &clientManager{
	baseClientMap:   map[string]*Client{},
	customClientMap: map[string]any{},
	lock:            sync.Mutex{},
}
