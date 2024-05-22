package client

import (
	"net/http"
)

type ReqContext struct {
	Headers     http.Header
	Params      map[string]string
	ContentType string
	Body        []byte
	HttpVersion string
}

type Service struct {
	Name      string   // 对应APPID
	Timeout   int      // 接口调用超时时间
	Headers   []string // 链路透传  Tracing
	Instances []string // 服务所对应的实例列表， 根据group, tag和zone进行筛选
	RateLimit int      // 限速配置
	Client    *http.Client
}

func (s *Service) choose() string {
	return ""
}

func (s *Service) Get(uri string, ctx ...ReqContext) {
	//req := http.NewRequest()
	//s.Client.Do(req)
}

func (s *Service) Post(uri string, ctx ...ReqContext) {
}

func (s *Service) Put(uri string, ctx ...ReqContext) {
}

func (s *Service) Delete(uri string, ctx ...ReqContext) {
}

func (s *Service) Head(uri string) {
}
