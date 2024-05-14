package conf

import (
	"encoding/json"
	"github.com/winjeg/go-commons/log"
)

type EventType byte

const (
	ConfigChange  = EventType(1)
	InfoChange    = EventType(2)
	SvcInfoChange = EventType(3)
)

const (
	TypeAdd    = "add"
	TypeRemove = "remove"
	TypeChange = "change"
)

type ConfigChangeEvent struct {
	Namespace string `json:"namespace,omitempty"`
	Key       string `json:"key,omitempty"`
	Type      string `json:"type,omitempty"`
	Current   string `json:"current,omitempty"`
	Before    string `json:"before,omitempty"`
}

type ConfigListener struct {
	Context *ConfigContext
}

func (l *ConfigListener) OnChange(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	eventType := data[len(data)-1]
	if eventType != byte(ConfigChange) {
		return nil
	}
	payloadData := data[:len(data)-1]
	cfgEvent := new(ConfigChangeEvent)
	if err := json.Unmarshal(payloadData, cfgEvent); err != nil {
		log.GetLogger(nil).Errorln("error reading event: ", err.Error())
		return err
	}

	switch cfgEvent.Type {
	case TypeChange, TypeAdd:
		l.Context.Lock.Lock()
		l.Context.Data[cfgEvent.Namespace][cfgEvent.Key] = cfgEvent.Current
		l.Context.Lock.Unlock()
	case TypeRemove:
		l.Context.Lock.Lock()
		delete(l.Context.Data[cfgEvent.Namespace], cfgEvent.Key)
		l.Context.Lock.Unlock()
	}
	return nil
}

func (l *ConfigListener) RetryTimes() int {
	return 1
}

func DefaultConfigListener(ctx *ConfigContext) *ConfigListener {
	return &ConfigListener{Context: ctx}
}
