package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func RpcRequest(rpcUrl, path, method string) ([]byte, http.Header, error) {
	p, err := url.JoinPath(rpcUrl, path)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest(method, p, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("RpcRequest: failed to fetch data")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("RpcRequest: failed to fetch data, statusCode: %v, error: %s", resp.StatusCode, err.Error())
		}
		return nil, nil, fmt.Errorf("RpcRequest: failed to fetch data, statusCode: %v, error: %s", resp.StatusCode, string(msgData))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return data, resp.Header, nil
}
