package httputil

import (
	"encoding/json"
	"net/http"

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
}

func NewRequest(method, url string, payload interface{}) *MyRequest {
	return &MyRequest{
		method:  method,
		url:     url,
		payLoad: payload,
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
	return post(req.url, req.payLoad)
}
