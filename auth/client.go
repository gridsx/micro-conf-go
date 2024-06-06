package auth

///  此client是为了给其他go客户端使用
/// 用来调用使用了openAPI签名校验规则的

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// 30 seconds timeout
	httpClient = &http.Client{Timeout: time.Second * 30}
)

type Client struct {
	Key    string
	Keeper SecretKeeper
}

func DefaultClient(key, sec string) *Client {
	return &Client{
		Key: key,
		Keeper: DefaultProvider{
			AppKey:    key,
			AppSecret: sec,
		},
	}
}

func (c *Client) Get(uri string, headers ...http.Header) (string, error) {
	return c.requestWithHeader(http.MethodGet, uri, nil, headers...)
}

func (c *Client) Post(uri string, body io.Reader, headers ...http.Header) (string, error) {
	return c.requestWithHeader(http.MethodPost, uri, body, headers...)
}

func (c *Client) Delete(uri string, body io.Reader, headers ...http.Header) (string, error) {
	return c.requestWithHeader(http.MethodDelete, uri, body, headers...)
}

func (c *Client) Put(uri string, body io.Reader, headers ...http.Header) (string, error) {
	return c.requestWithHeader(http.MethodPut, uri, body, headers...)
}

func (c *Client) requestWithHeader(method, url string, body io.Reader, headers ...http.Header) (string, error) {
	// 请求构造
	var header http.Header
	if len(headers) > 0 {
		header = headers[0]
	}
	secret, _ := c.Keeper.GetSecret(c.Key)
	url = BuildUrlParams(url, c.Key, secret)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	req.Header = header
	req.Close = true
	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return string(data), errors.New(fmt.Sprintf("Error with not correct status code %s", resp.Status))
	}
	data, err := io.ReadAll(resp.Body)
	defer safeClose(resp)
	return string(data), err
}

func BuildUrlParams(url, key, secret string) string {
	pairs := make([]KvPair, 0, 10)
	// add all params
	timeMillis := time.Now().UnixNano() / int64(time.Millisecond)
	paramPairs := getPairsFromUrl(url)
	pairs = append(pairs, paramPairs...)
	pairs = append(pairs, KvPair{
		Key:   timeParam,
		Value: fmt.Sprintf("%d", timeMillis),
	})
	pairs = append(pairs, KvPair{
		Key:   appKey,
		Value: key,
	})
	content := buildParams(pairs)
	signResult := Sign(content, secret)
	if len(getPairsFromUrl(url)) == 0 {
		url += fmt.Sprintf("?%s=%s&%s=%d&%s=%s", signParam, signResult, timeParam, timeMillis, appKey, key)
	} else {
		url += fmt.Sprintf("&%s=%s&%s=%d&%s=%s", signParam, signResult, timeParam, timeMillis, appKey, key)
	}
	return url
}

func safeClose(resp *http.Response) {
	if resp != nil && !resp.Close {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("error close response body:%v", err)
		}
	}
}

func getPairsFromUrl(rawUrl string) Pairs {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil
	}
	return getPairsFromMap(u.Query())
}

// get params and headers except the param sign
func getPairsFromMap(m map[string][]string) Pairs {
	pairs := make([]KvPair, 0, 10)
	for k, v := range m {
		if len(k) < 1 {
			continue
		}
		var val string
		for _, e := range v {
			val += e
		}
		if strings.EqualFold(k, signParam) {
			continue
		}
		p := KvPair{
			Key:   k,
			Value: val,
		}
		pairs = append(pairs, p)
	}
	return pairs
}
