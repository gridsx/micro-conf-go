# micro-conf-go

a go client for micro-conf

## 简介（introduction）

`micro-conf-go` is a go version client for `micro-conf`. It is designed to manage microservices and configs.
Using micro-conf you don't need a config center or a microservice register center. Everything is managed by micro-conf.

`micro-conf-go` 是一个go版本的 `micro-conf` 客户端， 而`micro-conf`是配置中心和注册中心的二合一， 有了它，你不再需要一个配置中心，
也不需要一个注册中心，
也不需要额外的任何其他部署，比如配置中心和注册中心所需要的 zk或者数据库等资源， 这些统统不需要。
目前go版本的客户端实现了配置加载和读取，配置变更监听， 服务注册， 心跳保持，客户端调用等功能， 也可以传递链路等。
目前这个版本属于一个比较简单的版本， 但是所需要的基本功能都是全的。


## 使用方法 （Getting started）

### 一、 获取依赖 (Get dependency)

```shell
go get  github.com/grdisx/micro-conf-go 
```

### 二、编写配置文件 (write config file for micro-conf)

配置文件格式如下，也可以是其他格式的， 只要符合 [定义格式](./mc/mc_entity.go) 中 `MicroConf` 格式定义即可

```yaml
# 服务 id 保证唯一即可
id: OrderService
# raft集群地址，用英文逗号分割，中间不能有空格, 推荐至少部署3台
meta-servers: 10.10.10.44:8000,10.10.10.32:8000,10.10.10.14:8000
# 服务注册的端口
port: 8080
# 服务对应的配置名称， 可以写多个用英文逗号分割， 也可以共享配置，需要指定 app.group.namespace 到sharedNamespace
namespaces: app.props
# 服务分组，默认default
group: default
# 服务与 micro-conf 注册配置中心的心跳间隔
heartbeat-timeout: 10
# 在管理页面新建应用时候分配的 应用对应的token
token: 8kA1W63KOSOXs6x9z8q40w9MXRjZJb9k
# 要注册的meta信息
meta:
  # zone 固定meta， 表示区域信息， 可用于客户端筛选zone调用
  zone: hz
  # tags 也是固定meta，表示标签， 可用于客户端筛选tag调用
  tags: beta,test
# 服务信息，是否开启服务注册，如果不开启默认注册的实例状态是 DISABLED， 开启则为 UP
service:
  enabled: true

# 依赖的客户端列表
clients:
  - name: OrderService
    group: default
    tags: test
    zone: hz
    timeout: 10
```

### 三、使用方法 （Getting ready）

示例如下，您需要先准备好 micro-conf 的配置对象 `mc.MicroConf` ， 可以从任意地方读取。
读取后， 可以直接使用 `mc.NewClient(conf)` 创建客户端， 然后启动客户端即可， 使用方法非常简单。

micro-conf 存在的意义就是为了合并配置中心和服务注册中心，减少对应的服务器成本和维护成本。因此他主要包括如下两方面的功能

1. 配置中心
    - 配置管理和获取
    - 配置变化监听
2. 注册中心
    - 把自身注册为可供其他应用调用的微服务
    - 生成可以调用其他微服务的客户端

```go
package demo

import (
	"fmt"
	"github.com/grdisx/micro-conf-go/mc"
)

func demo() {
	// conf 配置构造这一步需要自己去实现，
	conf := new(mc.MicroConf)

	// 创建client 并且启动， 正常来讲，需要在你的项目启动前启动
	client1 := mc.NewClient(conf)
	client1.Start()

	// 添加配置变更监听器
	client1.AddListener(xxxx)

	// ---------------- 配置中心 部分的功能用法 ------------------------
	// 从配置中心获取 key为 app.id对应的 int 值
	appId, _ := client1.GetInt("app.id")
	fmt.Println(appId)
	// 从配置中心获取 key为 userIds 对应的 int 值列表
	ids := client1.GetIntList("userIds")
	fmt.Println(ids)
	// 从配置中心获取一个map，key为 demo
	m := client1.GetAnyMap("demo")
	fmt.Println(m)
	// 从配置中心获取一个对象，目前是用json反序列化的	
	demo := new(DemoStruct)

	// -----  微服务调用方面的支持 ------------------------------------- 
	// 调用其他注册服务的客户端构造，(本例为自己调用自己)
	// 目前的负载均衡策略是 RoundRobin 
	orderService := client1.GetClient("OrderService")
	resp, err := orderService.Get(context.Background(), "/api/ping")
	orderService.Post(context.Background(), "/api/ping")
	orderService.Put(context.Background(), "/api/ping")
}
```

> 注: 本包支持，多个服务的注册，只需要多份配置即可， 配置上不能有冲突，比如端口不能冲突， app名称不能冲突


## 贡献 （how to contribute)

由于目前此项目仅有我一个人在维护和跟进， 如果需要新的功能，可能需要比较长的时间开发。 如果您对这个项目有兴趣，欢迎[发邮件](winjeg@qq.com)联系我

Contact me via my [email](winjeg@qq.com)

## 附录：与配置中心相连的API

1. 服务器启动的时候注册实例的API
2. 服务启动的时候Load配置的API
3. 监听WS消息的API
4. 初始化客户端的时候，获取对应服务实例列表的API