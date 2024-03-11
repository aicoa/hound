/*
 * @Author: aicoa
 * @Date: 2024-02-05 23:08:13
 * @Last Modified by: aicoa
 * @Last Modified time: 2024-02-20 00:35:13
 */
package client

import (
	"context"
	"crypto/tls"
	"hound/common"
	"io"
	"net/http"
	"time"
)

// 返回body信息
func NewHttpWithDefualtHeader(method, url string, client *http.Client) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultWebTimeout*time.Second)
	defer cancel()
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	if resp != nil && resp.StatusCode != 302 /*资源本来存在但是被更改位置*/ {
		defer resp.Body.Close()
		if aicoa, err := io.ReadAll(resp.Body); err == nil {
			return resp, aicoa, nil
		} else {
			return nil, nil, err
		}
	}
	return nil, nil, err
}

// 随页面跳转
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //考虑到网络环境因素，这里简单粗暴一点来防止https报错
		},
	}
}

// 不跟随页面跳转
func DefaultClientNoRedirect() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
