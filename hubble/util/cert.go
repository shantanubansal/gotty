package util

import (
	"crypto/tls"
	nethttp "net/http"
)

func GetHttpClientWithTls() *nethttp.Client {
	return &nethttp.Client{
		Transport: &nethttp.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
