/*
 * @Author: aicoa
 * @Date: 2024-02-19 22:29:54
 * @Last Modified by: aicoa
 * @Last Modified time: 2024-02-19 22:50:02
 */
package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"hound/common"
	"hound/common/logger"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// 选择代理模式，返回http.client

func SelectProxy(profile *common.Profiles) (client *http.Client) {
	if profile.Proxy.Mode == "HTTP" {
		urli := url.URL{}
		urlproxy, _ := urli.Parse(fmt.Sprintf("%v://%v:%v", profile.Proxy.Mode, profile.Proxy.Address, profile.Proxy.Port))
		client = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyURL(urlproxy),
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, //不用https
				TLSHandshakeTimeout: time.Second * 4,
			},
		}
	} else {
		auth := &proxy.Auth{
			User:     profile.Proxy.Username,
			Password: profile.Proxy.Password,
		}
		dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", profile.Proxy.Address, profile.Proxy.Port), auth, proxy.Direct)
		if err != nil {
			logger.Error(err)
		}
		httpTransport := &http.Transport{
			Dial:                dialer.Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: time.Second * 4,
		}
		client = &http.Client{
			Transport: httpTransport,
		}
		//设置sock5
		httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			con, err := dialer.Dial(network, addr)
			if err != nil {
				logger.Debug(err)
				return nil, err
			}
			return con, nil
		}
	}
	return client
}
