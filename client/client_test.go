package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinParams(t *testing.T) {
	assert.Equal(t, joinParams([]string{}), "")
	assert.Equal(t, joinParams([]string{"a"}), "")
	assert.Equal(t, joinParams([]string{"k", "v"}), "k=v")
	assert.Equal(t, joinParams([]string{"k", "v", "t"}), "k=v")
	assert.Equal(t, joinParams([]string{"k", "v", "t", "x"}), "k=v&t=x")
}

// 示例 Response 类型
type Response1 struct {
	BaseResponse
	Data Data1 `json:"data"` // 具体的数据字段
}

type Data1 struct {
	Field1 string `json:"field1"` // 数据字段
}

type Response2 struct {
	BaseResponse
	Data Data2 `json:"data"` // 具体的数据字段
}

type Data2 struct {
	Field2 int `json:"field2"` // 数据字段
}

// ExampleClient 是一个示例客户端结构体，包含多个方法.
type ExampleClient struct {
	Method1 func(name string, id int, body []byte) Response1         `method:"POST" uri:"/api/method1" content-type:"application/x-www-form-urlencoded"`
	Method2 func(id int, name string, body interface{}) Response2    `method:"POST" uri:"/api/method2"`
	Method3 func(data string, body map[string]interface{}) Response1 `method:"PUT" uri:"/api/method3"`
	Method4 func(id int, body []byte) Response2                      `method:"DELETE" uri:"/api/method4"`
}

func TestMakeClient(t *testing.T) {
	// 创建 ExampleClient 实例.
	client := &ExampleClient{}
	cm := &clientMaker{cli: &Client{}}
	// 为所有方法生成实现.
	cm.generateMethods(client)

	// 调用新生成的方法.
	if client.Method1 != nil {
		response1 := client.Method1("JohnDoe", 123, []byte(`name=JohnDoe&age=30`))
		fmt.Printf("Method1 Response: %+v\n", response1)
	}

	if client.Method2 != nil {
		type BodyStruct struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		}
		response2 := client.Method2(456, "JaneDoe", BodyStruct{Name: "JaneDoe", ID: 456})
		fmt.Printf("Method2 Response: %+v\n", response2)
	}

	if client.Method3 != nil {
		response3 := client.Method3("SomeData", map[string]interface{}{"key": "value"})
		fmt.Printf("Method3 Response: %+v\n", response3)
	}

	if client.Method4 != nil {
		response4 := client.Method4(789, []byte(`{"key":"value"}`))
		fmt.Printf("Method4 Response: %+v\n", response4)
	}
}
