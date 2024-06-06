package conf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grdisx/micro-conf-go/auth"
	"github.com/grdisx/micro-conf-go/utils"
	"net/http"
	"strings"
	"sync"
)

type NamespaceClientRequest struct {
	AppId            string   `json:"appId,omitempty"`
	Group            string   `json:"group,omitempty"`
	Namespaces       []string `json:"namespaces,omitempty"`
	SharedNamespaces []string `json:"sharedNamespaces,omitempty"`
}

type ConfigResult struct {
	Code string                       `json:"code,omitempty"`
	Msg  string                       `json:"msg,omitempty"`
	Data map[string]map[string]string `json:"data,omitempty"`
}

type ConfigContext struct {
	Data map[string]map[string]string
	Lock sync.Mutex
}

// Load is called once when application start up
func (c *ConfigContext) Load(req *NamespaceClientRequest, metaServers, token string) error {
	// choose addr, then add some
	server := utils.ChooseAddr(metaServers)
	addr := fmt.Sprintf("http://%s/api/cfg/app", server)
	d, err := json.Marshal(req)
	if err != nil {
		return err
	}
	cli := auth.DefaultClient(req.AppId, token)
	resp, respErr := cli.Post(addr, bytes.NewReader(d), http.Header{"Content-Type": []string{"application/json"}})
	if respErr != nil {
		return respErr
	}
	cfgResult := new(ConfigResult)
	if err := json.Unmarshal([]byte(resp), cfgResult); err != nil {
		return err
	}
	if strings.EqualFold(cfgResult.Code, "0") {
		c.Data = cfgResult.Data
		return nil
	}
	return errors.New(cfgResult.Msg)
}

func (c *ConfigContext) get(key string) (string, error) {
	for _, m := range c.Data {
		if v, ok := m[key]; ok {
			return v, nil
		}
	}
	return "", errors.New("not found")
}
