package httputil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"wstester/pkg/log"

	"github.com/pkg/errors"
)

func post(urlStr string, params interface{}) ([]byte, error) {
	var buffer []byte

	log.Infof("params: %+v", params)
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "marshal params errors")
	}

	content := bytes.NewReader(paramsBytes)

	//log.Infof("post urlStr: %s; content: %v", urlStr, content)
	//proxyUrl, err := url.Parse("http://10.1.13.121:50061")
	//if err != nil {
	//	log.Errorf("parse url error: %s", err.Error())
	//}
	//log.Infof("parse url successfully")
	//t := &http.Transport{
	//	MaxIdleConns:    10,
	//	MaxConnsPerHost: 10,
	//	IdleConnTimeout: 10 * time.Second,
	//	Proxy:           http.ProxyURL(proxyUrl),
	//}
	//cli := http.Client{Timeout: 10 * time.Second, Transport: t}
	cli := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPost, urlStr, content)
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
