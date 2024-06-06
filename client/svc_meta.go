package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grdisx/micro-conf-go/auth"
	"net/http"
	"strings"

	"github.com/winjeg/go-commons/log"
	"github.com/winjeg/go-commons/str"
)

const stateUP = "UP"

var logger = log.GetLogger(nil)

type SvcRet struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg,omitempty"`
	Data []ServiceInfo `json:"data,omitempty"`
}

type ServiceInfo struct {
	App              string                 `json:"app"`
	Group            string                 `json:"group"`
	IP               string                 `json:"ip,omitempty"`
	Port             int                    `json:"port,omitempty"`
	State            string                 `json:"state,omitempty"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
	HeartbeatTimeout int                    `json:"timeout,omitempty"`
}

func getServiceMeta(instance, app, group, token string) []ServiceInfo {
	url := fmt.Sprintf("http://%s/api/svc/instances", instance)
	info := ServiceInfo{App: app, Group: group}
	d, err := json.Marshal(info)
	if err != nil {
		logger.Errorf("getServiceMeta - json error: %s\n", err.Error())
		return nil
	}
	cli := auth.DefaultClient(app, token)
	resp, err := cli.Post(url, bytes.NewReader(d), http.Header{"Content-Type": []string{"application/json"}})
	if err != nil {
		logger.Errorf("getServiceMeta - remote call error:  %+v\n", err)
		return nil
	}

	svcRet := new(SvcRet)
	if err := json.Unmarshal([]byte(resp), svcRet); err != nil {
		logger.Errorf("getServiceMeta - err unmarshal resp data: %s\n", resp)
		return nil
	}

	if svcRet.Code != "0" {
		logger.Errorf("getServiceMeta - err unmarshal resp data: %s\n", resp)
		return nil
	}
	return svcRet.Data
}

// 只取过滤后的，以及状态为UP的服务
// zone 只能精准过滤
// tags可以交叉过滤，只要服务器包含客户端的任一tag即可
func processServiceInstance(metas []ServiceInfo, filters map[string]string) []string {
	instances := make([]string, 0, len(metas))
	for _, v := range metas {
		if len(filters) > 0 {
			filterSuccess := false
			if tags, ok := filters[tagFiler]; ok {
				if len(v.Meta) == 0 {
					continue
				}
				if metaTagStr, ok := v.Meta[tagFiler]; ok {
					if metaTagStr == nil {
						continue
					}
					if mts, ok := metaTagStr.(string); ok {
						metaTags := strings.Split(mts, ",")
						for _, mt := range metaTags {
							clientTags := strings.Split(tags, ",")
							if str.Contains(clientTags, mt) {
								// 包含有tag
								filterSuccess = true
								break
							}
						}
					}
				}
			} else {
				// 没有tag filter
				filterSuccess = true
			}
			if zone, ok := filters[zoneFilter]; ok {
				if len(v.Meta) == 0 {
					continue
				}
				if metaZoneStr, ok := v.Meta[zoneFilter]; ok {
					if metaZone, ok := metaZoneStr.(string); ok {
						if !strings.EqualFold(metaZone, zone) {
							filterSuccess = false
						}
					} else {
						filterSuccess = false
					}
				} else {
					filterSuccess = false
				}
			}

			if filterSuccess && v.State == stateUP {
				inst := fmt.Sprintf("%s:%d", v.IP, v.Port)
				instances = append(instances, inst)
			}
		} else {
			if v.State == stateUP {
				inst := fmt.Sprintf("%s:%d", v.IP, v.Port)
				instances = append(instances, inst)
			}
		}
	}
	return instances
}
