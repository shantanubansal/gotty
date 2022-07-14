package util

import (
	"crypto/tls"
	nethttp "net/http"
)

func GetTlsConfig(certificatePath, keyPath string, isInsecure bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certificatePath, keyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: isInsecure,
	}, nil
}

func GetCertWithPath(certificate, key string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	return cert, nil
}
func GetTlsWithPath(certificate, key string, insecureSkipVerify bool) *tls.Config {
	cert, err := GetCertWithPath(certificate, key)
	if err != nil {
		return nil
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: insecureSkipVerify,
	}
}
func GetHttpClientWithTls(certificatePath, keyPath string, insecureSkipVerify bool) *nethttp.Client {
	return &nethttp.Client{
		Transport: &nethttp.Transport{
			TLSClientConfig: GetTlsWithPath(certificatePath, keyPath, insecureSkipVerify),
		},
	}
}
