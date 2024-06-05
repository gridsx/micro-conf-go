package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/grdisx/micro-conf-go/utils"
)

const (
	defaultTimeout = 30
	tagFiler       = "tags"
	zoneFilter     = "zone"
)

type Client struct {
	Name          string            // 对应APPID
	Group         string            // 对应groupId
	Headers       map[string]string // 链路透传  Tracing
	Instances     []string          // 服务所对应的实例列表， 根据group, tag和zone进行筛选
	Filters       map[string]string // 过滤条件， 只有符合过滤条件的，才会被放入Instances
	RefreshPeriod int               // 配置刷新间隔
	Timeout       int               // 接口调用超时时间，单位秒
	RateLimit     int               // 限速配置
	Client        *http.Client
	instNum       atomic.Int32
}

func (s *Client) init(metaServers string) error {
	if s.Timeout <= 0 {
		s.Timeout = defaultTimeout
	}
	client := &http.Client{Timeout: time.Duration(s.Timeout) * time.Second}
	s.Client = client
	s.getInstances(metaServers)
	if len(s.Instances) == 0 {
		logger.Errorf("error getting meta info for svc:%s, group:%s\n", s.Name, s.Group)
		return errors.New(fmt.Sprintf("no instances for svc: %s, group:%s", s.Name, s.Group))
	}
	go func() {
		ticker := time.Tick(time.Duration(s.RefreshPeriod) * time.Second)
		for {
			select {
			case <-ticker:
				s.getInstances(metaServers)
			}
		}
	}()
	return nil
}

func (s *Client) getInstances(metaServers string) {
	addr := utils.ChooseAddr(metaServers)
	metas := getServiceMeta(addr, s.Name, s.Group, s.Client)
	if len(metas) == 0 {
		logger.Errorf("error getting meta info for svc:%s, group:%s\n", s.Name, s.Group)
		return
	}
	instances := processServiceInstance(metas, s.Filters)
	// 中间有一次get不到，不要覆盖
	if len(instances) > 0 {
		s.Instances = instances
	} else {
		logger.Errorf("error getting meta info for svc:%s, group:%s\n", s.Name, s.Group)
	}
}

// 轮询抽取
func (s *Client) choose() (string, error) {
	if len(s.Instances) == 0 {
		return "", errors.New("")
	}
	if int(s.instNum.Add(1)) >= len(s.Instances) {
		s.instNum.Store(0)
	}
	return s.Instances[int(s.instNum.Load())], nil
}

func (s *Client) Get(ctx context.Context, uri string) (*http.Response, error) {
	return s.request(ctx, http.MethodPut, uri, nil, []string{})
}

func (s *Client) GetWithParams(ctx context.Context, uri string, params ...string) (*http.Response, error) {
	return s.request(ctx, http.MethodPut, uri, nil, params)
}

func (s *Client) Post(ctx context.Context, uri string, body io.Reader, params ...string) (*http.Response, error) {
	return s.request(ctx, http.MethodPost, uri, body, params)
}

func (s *Client) Put(ctx context.Context, uri string, body io.Reader, params ...string) (*http.Response, error) {
	return s.request(ctx, http.MethodPut, uri, body, params)
}

func (s *Client) Delete(ctx context.Context, uri string, body io.Reader, params ...string) (*http.Response, error) {
	return s.request(ctx, http.MethodDelete, uri, body, params)
}

func (s *Client) Head(uri string) (*http.Response, error) {
	return s.Client.Head(uri)
}

func (s *Client) request(ctx context.Context, method, uri string, body io.Reader, params []string) (*http.Response, error) {
	url := uri
	if len(params) > 1 {
		url = url + "?" + joinParams(params)
	}
	addr, err := s.choose()
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("http://%s%s", addr, url)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = s.headers(ctx)
	return s.Client.Do(req)
}

func (s *Client) headers(ctx context.Context) http.Header {
	v := ctx.Value("headers")
	if len(s.Headers) == 0 && v == nil {
		return nil
	}
	headers := http.Header{}
	if s.Headers != nil {
		for k, v := range s.Headers {
			headers[k] = []string{v}
		}
	}
	if v != nil {
		if header, ok := v.(http.Header); ok {
			for k, v := range header {
				headers[k] = v
			}
			return headers
		}
	}
	return headers
}

func joinParams(params []string) string {
	if len(params) > 1 {
		if len(params)%2 != 0 {
			params = params[:len(params)-1]
		}
	} else {
		return ""
	}
	pairs := make([]string, 0, len(params)/2)
	for i := 0; i < len(params); i += 2 {
		p := fmt.Sprintf("%s=%s", params[i], params[i+1])
		pairs = append(pairs, p)
	}
	return strings.Join(pairs, "&")
}
