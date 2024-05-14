package conn

import (
	"strings"
)

type HeartBeat struct {
	AppId      string            `json:"appId,omitempty"`
	Group      string            `json:"group,omitempty"`
	IP         string            `json:"ip,omitempty"`
	Port       int               `json:"port,omitempty"`
	EnableSvc  bool              `json:"enableSvc"`
	EnableCfg  bool              `json:"enableCfg"`
	Meta       map[string]string `json:"meta,omitempty"`
	Namespaces []string          `json:"namespaces"` //此实例监听了哪些 配置文件， 如果是自己的话，取自己的appId和group， 如果是他人的，则取他人的
	Timeout    int64             `json:"timeout,omitempty"`
}

type RegAppPayload struct {
	AppId            string            `json:"appId,omitempty"`
	Group            string            `json:"group,omitempty"`
	IP               string            `json:"ip,omitempty"`
	Port             int               `json:"port,omitempty"`
	State            string            `json:"state,omitempty"`
	Meta             map[string]string `json:"meta,omitempty"`
	Namespaces       []string          `json:"namespaces,omitempty"`
	SharedNamespaces []string          `json:"sharedNamespaces,omitempty"`
}

type HttpResult struct {
	Code string      `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type EventType string

func (t *EventType) Byte() byte {
	if strings.EqualFold(string(*t), string(ConfigChange)) {
		return 1
	}
	if strings.EqualFold(string(*t), string(InfoChange)) {
		return 2
	}
	if strings.EqualFold(string(*t), string(SvcInfoChange)) {
		return 3
	}
	return 0
}

const (
	ConfigChange  = EventType("cfg")
	InfoChange    = EventType("info")
	SvcInfoChange = EventType("svc")
)

type AppEvent struct {
	Type    EventType   `json:"type,omitempty"`
	Content interface{} `json:"content,omitempty"`
}
