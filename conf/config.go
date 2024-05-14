package conf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
func (c *ConfigContext) Load(req *NamespaceClientRequest, metaServers string) error {
	// choose addr, then add some
	server := chooseAddr(metaServers)
	addr := fmt.Sprintf("http://%s/api/cfg/app", server)
	d, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, respErr := http.Post(addr, "application/json", bytes.NewReader(d))
	if respErr != nil {
		return respErr
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("error response code")
	}
	respData, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	cfgResult := new(ConfigResult)
	if err := json.Unmarshal(respData, cfgResult); err != nil {
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

func chooseAddr(metaServers string) string {
	addrArr := strings.Split(metaServers, ",")
	idx := rand.Intn(len(addrArr))
	return addrArr[idx]
}
