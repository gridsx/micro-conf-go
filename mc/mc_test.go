package mc

import (
	"context"
	"fmt"
	"github.com/grdisx/micro-conf-go/client"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const demoConf = `
id: OrderService
meta-servers: 10.10.10.44:8000
port: 8080
namespaces: app.props
group: default
heartbeat-timeout: 10
token: 03120382109Ha
meta:
  zone: hz
  tags: beta,test
service:
  enabled: true
clients:
  - name: OrderService
    group: default
    tags: test
    zone: hz
    timeout: 10
`

func TestStart(t *testing.T) {
	conf := new(MicroConf)
	err := yaml.Unmarshal([]byte(demoConf), conf)
	if err != nil {
		t.FailNow()
		return
	}
	client1 := NewClient(conf)
	client1.Start()
	conf.Port = 8003
	client2 := NewClient(conf)
	client2.Start()
	select {}
}

type DemoStruct struct {
	Name string `json:"name,omitempty"`
	Id   int    `json:"id,omitempty"`
}

func TestConfig(t *testing.T) {
	conf := new(MicroConf)
	err := yaml.Unmarshal([]byte(demoConf), conf)
	if err != nil {
		t.FailNow()
		return
	}
	client1 := NewClient(conf)
	client1.Start()
	_, intErr := client1.GetInt("app.id")
	assert.Nil(t, intErr)

	ids := client1.GetIntList("userIds")
	assert.True(t, len(ids) > 0)

	m := client1.GetAnyMap("demo")
	fmt.Println(m)

	str := client1.GetStringDefault("abc.123", "aaa")
	assert.Equal(t, str, "aaa")

	demo := new(DemoStruct)

	ticker := time.Tick(time.Second * 5)
	orderService := client.GetClient("OrderService")
	count := 3
	for {
		select {
		case <-ticker:
			if count <= 0 {
				return
			}
			client1.GetObject("demo", demo)
			fmt.Println(demo.Name, demo.Id)
			if resp, err := orderService.Get(context.Background(), "/api/ping"); err != nil {
				fmt.Println(err.Error(), resp)
			}
			count--
		}
	}
}
