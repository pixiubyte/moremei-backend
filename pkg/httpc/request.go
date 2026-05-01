package httpc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

func Request(method, url string, headers map[string]string, data []byte) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logx.Errorf("request err:%v", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
