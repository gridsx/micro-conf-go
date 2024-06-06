package mc

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/grdisx/micro-conf-go/client"
	"github.com/grdisx/micro-conf-go/conf"
	"github.com/grdisx/micro-conf-go/conn"
)

const (
	defaultSize = 8
)

type MicroClient struct {
	Cfg       *MicroConf
	SocketMgr *conn.WebsocketMgr
	Listeners []ChangeListener
	context   *conf.ConfigContext
	clientMap map[string]*client.Client
}

func (m *MicroClient) GetClient(serviceName string) *client.Client {
	if m.clientMap == nil {
		return nil
	}
	return m.clientMap[serviceName]
}

func (m *MicroClient) AddListener(l ChangeListener) {
	m.Listeners = append(m.Listeners, l)
}

func (m *MicroClient) Start() {
	c := m.Cfg
	regErr := m.registerApp()
	if regErr != nil {
		panic(regErr)
	}

	// 配置中心相关， 开机load并加监听
	m.context = &conf.ConfigContext{Data: map[string]map[string]string{}, Lock: sync.Mutex{}}
	if m.Cfg.CfgEnabled() {
		m.context.Load(m.Cfg.NamespaceReq(), m.Cfg.MetaServers, m.Cfg.Token)
		m.Listeners = append(m.Listeners, conf.DefaultConfigListener(m.context))
	}

	// 如果开启了客户端，需要构造刷新
	if len(m.Cfg.Clients) > 0 {
		m.initClients()
	}

	mgr := conn.InitSocketMgr(c.MetaServers, c.Id, c.Token, c.Port)
	m.SocketMgr = mgr
	mgr.Start()
	m.sendHeartbeat()
	m.acceptConfigChange()

}

func (m *MicroClient) registerApp() error {
	c := m.Cfg
	if !c.Valid() {
		return errors.New("invalid config")
	}
	if err := conn.RegApp(c.ToRegAppPayload(), c.MetaServers, c.Token); err != nil {
		return err
	}
	return nil
}

func (m *MicroClient) sendHeartbeat() {
	c := m.Cfg
	heartbeat := c.ToHeartBeat()
	d, _ := json.Marshal(heartbeat)
	ticker := time.Tick(time.Second * time.Duration(heartbeat.Timeout/2))
	go func() {
		for {
			select {
			case <-ticker:
				m.SocketMgr.Send(d)
			}
		}
	}()
}

func (m *MicroClient) acceptConfigChange() {
	go func() {
		for {
			select {
			case d := <-m.SocketMgr.Recv():
				for _, l := range m.Listeners {
					go func(d []byte, l ChangeListener) {
						err := l.OnChange(d)
						retryTimes := l.RetryTimes()
						for err != nil && retryTimes > 0 {
							retryTimes--
							err = l.OnChange(d)
						}
					}(d, l)
				}
			}
		}
	}()
}

// 获取服务列表， 拿出来status是UP的，然后进行筛选
func (m *MicroClient) initClients() {
	m.clientMap = make(map[string]*client.Client, defaultSize)
	for _, v := range m.Cfg.Clients {
		filters := make(map[string]string, defaultSize)
		filters["tags"] = v.Tags
		filters["zone"] = v.Zone
		cli := client.NewClient(v.Name, v.Group, m.Cfg.MetaServers, m.Cfg.Token,
			filters, v.Headers, v.Lease, v.Timeout, v.RateLimit)
		m.clientMap[v.Name] = cli
	}
}

func NewClient(c *MicroConf) *MicroClient {
	return &MicroClient{Cfg: c, Listeners: []ChangeListener{}}
}
