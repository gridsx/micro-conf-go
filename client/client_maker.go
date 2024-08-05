package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const defaultSize = 4

// BaseResponse 定义一个通用的 Response 结构体基类
type BaseResponse struct {
	Status  int    `json:"status"`  // HTTP状态码
	Success bool   `json:"success"` // 是否成功
	Msg     string `json:"msg"`
}

type clientMaker struct {
	cli *Client
	// TODO monitoring support
}

// generateMethods 为给定结构体中的所有方法生成默认实现.
func (cm *clientMaker) generateMethods(obj interface{}) {
	val := reflect.ValueOf(obj)
	typ := val.Type()

	// 确保我们有一个指向结构体的指针.
	if val.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		// 获取实际的结构体值.
		val = val.Elem()
		typ = typ.Elem()
	}

	// 遍历结构体的所有字段.
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// 检查字段是否是一个方法.
		if field.Type.Kind() == reflect.Func && field.PkgPath == "" { // 确保这是一个公开的方法
			// 获取方法的值.
			methodVal := val.FieldByName(field.Name)

			// 检查方法是否已被实现.
			if methodVal.IsNil() {
				// 检查方法是否可寻址.
				if val.CanAddr() {
					// 生成方法实现.
					newFunc := reflect.MakeFunc(field.Type, cm.generateMethodImplementation(field))
					// 设置方法的值.
					val.Field(i).Set(newFunc)
				}
			}
		}
	}
}

// generateMethodImplementation 为特定方法生成实现.
func (cm *clientMaker) generateMethodImplementation(field reflect.StructField) func(args []reflect.Value) []reflect.Value {

	// 根据方法的标签信息构建请求.
	methodTag := "method"
	uriTag := "uri"
	contentTypeTag := "content-type"

	method, _ := field.Tag.Lookup(methodTag)
	uri, _ := field.Tag.Lookup(uriTag)
	contentType, _ := field.Tag.Lookup(contentTypeTag)

	// 默认的 Content-Type
	if contentType == "" {
		contentType = "application/json"
	}

	return func(args []reflect.Value) []reflect.Value {
		// 构建请求.
		queryParams := url.Values{}
		bodyContent := []byte{}
		headerMap := make(map[string]string, defaultSize)

		// 处理方法参数.
		for i, arg := range args {
			// 获取命名参数的名字.
			paramName := field.Type.In(i).Name()
			if i == 0 && arg.Kind() == reflect.Map {
				firstArgIsHeader := true
				if len(args) == 1 {
					// 如果只有一个参数， 且命名为headers， 那么就认为是header map， 否则认为是body
					if !strings.EqualFold(paramName, "headers") {
						firstArgIsHeader = false
					}
				}
				if m := arg.Interface(); m != nil && firstArgIsHeader {
					if v, ok := m.(map[string]string); ok {
						headerMap = v
					}
				}
			}

			// 添加到查询参数中.
			switch arg.Kind() {
			case reflect.Slice:
				if i == len(args)-1 && arg.Type().Elem().Kind() == reflect.Uint8 {
					bodyContent = arg.Bytes()
				} else if i == len(args)-1 && arg.Type().Elem().Kind() != reflect.Uint8 {
					bodyBytes, err := json.Marshal(arg.Interface())
					if err != nil {
						fmt.Printf("Error marshaling to JSON: %v\n", err)
						return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
					}
					bodyContent = bodyBytes
				}
			case reflect.Map:
				if i == len(args)-1 {
					bodyBytes, err := json.Marshal(arg.Interface())
					if err != nil {
						fmt.Printf("Error marshaling to JSON: %v\n", err)
						return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
					}
					bodyContent = bodyBytes
				}
			// 可以添加更多类型的处理.
			default:
				// 防止 nil 引用.
				queryParams.Add(paramName, arg.String())
			}
		}

		baseURL, instErr := cm.cli.choose()
		if instErr != nil {
			respType := field.Type.Out(0)
			response := reflect.New(respType).Elem()
			response.FieldByName("Status").SetInt(int64(-1))
			response.FieldByName("Success").SetBool(false)
			response.FieldByName("Msg").SetString("no available instance")
			return []reflect.Value{response}
		}

		// 构造完整 URL.
		var fullURL string
		if len(bodyContent) > 0 {
			// 如果有 body，则不需要 query params.
			fullURL = baseURL + uri
		} else {
			// 如果没有 body，则使用 query params.
			fullURL = baseURL + uri + "?" + queryParams.Encode()
		}

		// 发起 HTTP 请求.
		var req *http.Request
		var err error
		switch method {
		case "GET":
			req, err = http.NewRequest("GET", fullURL, nil)
		case "POST":
			req, err = http.NewRequest("POST", fullURL, bytes.NewReader(bodyContent))
		case "PUT":
			req, err = http.NewRequest("PUT", fullURL, bytes.NewReader(bodyContent))
		case "DELETE":
			req, err = http.NewRequest("DELETE", fullURL, nil)
		default:
			err = fmt.Errorf("unsupported HTTP method: %s", method)
		}

		if err != nil {
			fmt.Printf("Error creating Request: %v\n", err)
			return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
		}

		// 设置 Content-Type
		req.Header.Set("Content-Type", contentType)
		if len(headerMap) > 0 {
			for k, v := range headerMap {
				req.Header.Set(k, v)
			}
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error performing Request: %v\n", err)
			return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
		}
		defer resp.Body.Close()

		// 读取响应体.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
		}

		// 反序列化响应体到 data 字段.
		var data interface{}
		if resp.StatusCode < 400 {
			err = json.Unmarshal(respBody, &data)
			if err != nil {
				fmt.Printf("Error unmarshaling response: %v\n", err)
				return []reflect.Value{reflect.New(field.Type.Out(0)).Elem()}
			}
		}

		// 获取方法的返回类型.
		respType := field.Type.Out(0)
		response := reflect.New(respType).Elem()
		response.FieldByName("Status").SetInt(int64(resp.StatusCode))
		response.FieldByName("Success").SetBool(resp.StatusCode < 400)
		if resp.StatusCode < 400 {
			response.FieldByName("Msg").SetString("success")
		}

		// 获取 Data 字段的类型.
		if dft, ok := respType.FieldByName("Data"); ok {
			// 检查 Data 字段是否存在.
			if _, ok := response.Type().FieldByName("Data"); !ok {
				fmt.Println("Data field not found in the response type.")
				return []reflect.Value{response}
			}
			// 将 data 字段反序列化为正确的类型.
			dataValue := reflect.New(dft.Type).Elem()
			if err := json.Unmarshal(respBody, dataValue.Addr().Interface()); err == nil {
				// 设置 Data 字段.
				response.FieldByName("Data").Set(dataValue)
			} else {
				response.FieldByName("Success").SetBool(false)
				response.FieldByName("Msg").SetString(err.Error())
			}
		}
		// 返回响应结构体.
		return []reflect.Value{response}
	}
}
