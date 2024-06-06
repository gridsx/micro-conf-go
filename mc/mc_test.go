package mc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
token: 8kA1W63KOSOXs6x9z8q40w9MXRjZJb9k
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

	go serveHttp(conf.Port)

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
	orderService := client1.GetClient("OrderService")
	count := 3
	for {
		select {
		case <-ticker:
			if count <= 0 {
				return
			}
			client1.GetObject("demo", demo)
			fmt.Println(demo.Name, demo.Id)
			resp, err := orderService.Get(context.Background(), "/api/ping")
			if err != nil {
				fmt.Println(err.Error())
			} else {
				d, _ := io.ReadAll(resp.Body)
				fmt.Println(string(d))
			}
			count--
		}
	}
}

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func ping(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "pong")
}

func serveHttp(port int) {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/api/ping", ping)
	// 启动web服务，监听9090端口
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
