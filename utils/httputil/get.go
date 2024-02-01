package httputil

import (
	"bufio"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func get(urlStr string) ([]byte, error) {
	var buffer []byte

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

	//cli := &http.Client{Timeout: 30 * time.Second, Transport: t}
	cli := &http.Client{Timeout: 30 * time.Second}

	request, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request error")
	}

	resp, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "send request error")
	}

	defer resp.Body.Close()
	buffer = make([]byte, 1024)
	reader := bufio.NewReader(resp.Body)
	count, err := reader.Read(buffer)
	if err != nil {
		return nil, errors.Wrap(err, "read resp error")
	}

	return buffer[:count], nil
}
