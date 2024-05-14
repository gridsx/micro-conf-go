package client

type Service struct {
	Name      string   // 对应APPID
	Instances []string // 服务所对应的实例列表
	Timeout   int
}

func (s *Service) choose() string {
	return ""
}

func (s *Service) Get(uri string) {
}

func (s *Service) GetWithParams(uri string, params ...string) {
}

func (s *Service) Post(uri, contentType string, body []byte, params ...string) {
}

func (s *Service) PostJson(uri string, body []byte, params ...string) {
}

func (s *Service) Put(uri, contentType string, body []byte, params ...string) {
}

func (s *Service) PutJson(uri string, body []byte, params ...string) {
}

func (s *Service) Delete(uri, contentType string, body []byte, params ...string) {
}

func (s *Service) DeleteJson(uri string, body []byte, params ...string) {
}

func (s *Service) Head(uri string) {
}
