package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func HttpRequest(rpcUrl, path, method string, isApi bool) ([]byte, http.Header, error) {
	p, err := url.JoinPath(rpcUrl, path)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest(method, p, nil)
	if err != nil {
		return nil, nil, err
	}
	if isApi {
		req.Header.Set("Accept", "application/json")
	} else {
		req.Header.Set("Accept-Encoding", "gzip, deflate, br, identity")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("RpcRequest-> url: %s, error: %s", p, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, error: %s", p, resp.StatusCode, err.Error())
		}
		return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, body: %s", p, resp.StatusCode, string(msgData))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return data, resp.Header, nil
}
