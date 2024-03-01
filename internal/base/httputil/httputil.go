package httputil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"wstester/pkg/log"

	"github.com/pkg/errors"
)

type HttpUtil struct {
}

func NewHttpUtil() *HttpUtil {
	return new(HttpUtil)
}

type MyRequest struct {
	method   string
	url      string
	urlParam interface{}
	payLoad  interface{}
	isPublic bool
}

func NewRequest(method, url string, payload interface{}, isPublic bool) *MyRequest {
	return &MyRequest{
		method:   method,
		url:      url,
		payLoad:  payload,
		isPublic: isPublic,
	}
}

func (h HttpUtil) SendRequest(req *MyRequest, resp interface{}) error {
	var (
		body []byte
		err  error
	)

	switch req.method {
	case http.MethodGet:
		body, err = h.handleGet(req)
	case http.MethodPost:
		body, err = h.handlePost(req)
	default:
		return errors.New("unsupported method: " + req.method)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(body, resp)
}

func (h HttpUtil) handleGet(req *MyRequest) ([]byte, error) {
	return get(req.url)
}

func (h HttpUtil) handlePost(req *MyRequest) ([]byte, error) {
	// urlStr string, params interface{}
	var buffer []byte

	log.Infof("params: %+v", req.payLoad)
	paramsBytes, err := json.Marshal(req.payLoad)
	if err != nil {
		return nil, errors.Wrap(err, "marshal params errors")
	}

	content := bytes.NewReader(paramsBytes)

	// log.Infof("post urlStr: %s; content: %v", urlStr, content)
	cli := http.Client{Timeout: 10 * time.Second}
	if req.isPublic {
		proxyUrl, err := url.Parse("http://10.1.13.121:50061")
		if err != nil {
			log.Errorf("parse url error: %s", err.Error())
		}
		log.Infof("parse url successfully")
		cli.Transport = &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: 10 * time.Second,
			Proxy:           http.ProxyURL(proxyUrl),
		}
	}

	request, err := http.NewRequest(http.MethodPost, req.url, content)
	if err != nil {
		return nil, errors.Wrap(err, "new request error")
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "send post request error")
	}

	defer resp.Body.Close()
	buffer = make([]byte, 1024)
	reader := bufio.NewReader(resp.Body)
	count, err := reader.Read(buffer)
	if err != nil {
		return nil, errors.Wrap(err, "read resp error")
	}

	log.Infof("result: %s", string(buffer[:count]))
	return buffer[:count], nil
}
