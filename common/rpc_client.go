package common

import (
	"fmt"
	"io"
	"net/http"
)

func ApiRequest(rpcUrl, method string) ([]byte, http.Header, error) {
	maxRetries := 1
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest(method, rpcUrl, nil)
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
			// return nil, nil, fmt.Errorf("RpcRequest-> url: %s, error: %s", rpcUrl, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			msgData, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, error: %s", rpcUrl, resp.StatusCode, err.Error())
			}
			return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, body: %s", rpcUrl, resp.StatusCode, string(msgData))
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		return data, resp.Header, nil
	}
	return nil, nil, fmt.Errorf("RpcRequest: failed to fetch data after %d retries", maxRetries)
}

func HtmlRequest(rpcUrl string) ([]byte, http.Header, error) {
	maxRetries := 1
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", rpcUrl, nil)
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
			// return nil, nil, fmt.Errorf("RpcRequest-> url: %s, error: %s", rpcUrl, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			msgData, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, error: %s", rpcUrl, resp.StatusCode, err.Error())
			}
			return nil, nil, fmt.Errorf("RpcRequest-> url: %s, statusCode: %v, body: %s", rpcUrl, resp.StatusCode, string(msgData))
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		return data, resp.Header, nil
	}
	return nil, nil, fmt.Errorf("RpcRequest: failed to fetch data after %d retries", maxRetries)
}
