package mc

import (
	"github.com/grdisx/micro-conf-go/conf"
	"strings"

	"github.com/grdisx/micro-conf-go/conn"
	"github.com/grdisx/micro-conf-go/utils"
)

const (
	minHeartbeatTimeout     = 3
	defaultHeartBeatTimeout = 30
)

// ClientConfig 服务下线不会主动通知 ，因此正常下线要先把state设置为down超过调用方的lease后才可真正关闭服务
type ClientConfig struct {
	Name      string `json:"name,omitempty" yaml:"name"`           // 调用的服务名
	Group     string `json:"group,omitempty" yaml:"group"`         // 调用服务的分组， 不写默认为 default
	Lease     int64  `json:"lease,omitempty" yaml:"lease"`         // 本地缓存的服务的过期时间，即刷新时间， 推荐小于对应服务的心跳时间
	Tags      string `json:"tags,omitempty" yaml:"tags"`           // 筛选的tag， 只调用具备相应tag的节点， 不写则为全部
	Zones     string `json:"zones,omitempty" yaml:"zones"`         // 筛选的zone，只调用具备相应zone的节点， 不写则为全部
	RateLimit int    `json:"rateLimit,omitempty" yaml:"rateLimit"` //总限速， 超过限速会返回error， 小于0 或者不写均为不限速
}

type ServiceConfig struct {
	Enabled bool `json:"enabled,omitempty" yaml:"enabled"`
}

type MicroConf struct {
	Id               string            `json:"id,omitempty" yaml:"id"`
	MetaServers      string            `json:"metaServers,omitempty" yaml:"meta-servers"`
	Port             int               `json:"serverPort,omitempty" yaml:"port"`
	Group            string            `json:"group,omitempty" yaml:"group"`
	Token            string            `json:"token,omitempty" yaml:"token"`
	Meta             map[string]string `json:"meta,omitempty" yaml:"meta"`
	HeartbeatTimeout int64             `json:"heartbeatTimeout,omitempty" yaml:"heartbeat-timeout"`
	Namespaces       string            `json:"namespaces,omitempty" yaml:"namespaces"`
	SharedNamespaces string            `json:"sharedNamespaces,omitempty" yaml:"shared-namespaces"`
	Service          *ServiceConfig    `json:"service" yaml:"service"`
	Clients          []ClientConfig    `json:"clients,omitempty" yaml:"clients"`
}

func (m *MicroConf) Valid() bool {
	return len(strings.TrimSpace(m.Id)) > 0 && len(strings.TrimSpace(m.MetaServers)) > 0
}

func (m *MicroConf) extractNamespaces(str string) []string {
	namespaces := make([]string, 0, defaultSize)
	namespaceArr := strings.Split(str, ",")
	if len(namespaceArr) > 0 {
		for _, v := range namespaceArr {
			namespace := strings.TrimSpace(v)
			if !strings.EqualFold(namespace, "") {
				namespaces = append(namespaces, namespace)
			}
		}
	}
	return namespaces
}

func (m *MicroConf) ToRegAppPayload() *conn.RegAppPayload {

	group := m.Group
	if len(m.Group) == 0 {
		group = "default"
	}

	state := "DISABLED"

	if m.Service != nil && m.Service.Enabled {
		state = "UP"
	}

	namespaces := m.extractNamespaces(m.Namespaces)
	sharedNamespaces := m.extractNamespaces(m.SharedNamespaces)

	return &conn.RegAppPayload{
		AppId: m.Id,
		Group: group,
		IP:    utils.GetIP(),

		State: state,
		Port:  m.Port,
		Meta:  m.Meta,

		Namespaces:       namespaces,
		SharedNamespaces: sharedNamespaces,
	}

}

func (m *MicroConf) CfgEnabled() bool {
	namespaces := make([]string, 0, defaultSize)
	ns := m.extractNamespaces(m.Namespaces)
	ss := m.extractNamespaces(m.SharedNamespaces)
	if len(ns) > 0 {
		namespaces = append(namespaces, ns...)
	}
	if len(ss) > 0 {
		namespaces = append(namespaces, ss...)
	}
	return len(namespaces) > 0
}

func (m *MicroConf) ToHeartBeat() *conn.HeartBeat {

	group := m.Group
	if len(m.Group) == 0 {
		group = "default"
	}
	namespaces := make([]string, 0, defaultSize)
	ns := m.extractNamespaces(m.Namespaces)
	ss := m.extractNamespaces(m.SharedNamespaces)
	if len(ns) > 0 {
		namespaces = append(namespaces, ns...)
	}
	if len(ss) > 0 {
		namespaces = append(namespaces, ss...)
	}

	if m.HeartbeatTimeout < minHeartbeatTimeout {
		m.HeartbeatTimeout = defaultHeartBeatTimeout
	}

	return &conn.HeartBeat{
		AppId:      m.Id,
		Group:      group,
		IP:         utils.GetIP(),
		Port:       m.Port,
		EnableSvc:  m.Service != nil && m.Service.Enabled,
		EnableCfg:  len(namespaces) > 0,
		Meta:       m.Meta,
		Namespaces: namespaces,
		Timeout:    m.HeartbeatTimeout,
	}
}

func (m *MicroConf) NamespaceReq() *conf.NamespaceClientRequest {
	group := m.Group
	if len(m.Group) == 0 {
		group = "default"
	}
	return &conf.NamespaceClientRequest{
		AppId:            m.Id,
		Group:            group,
		Namespaces:       m.extractNamespaces(m.Namespaces),
		SharedNamespaces: m.extractNamespaces(m.SharedNamespaces),
	}
}
