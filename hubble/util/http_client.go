package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	basicAuth = "BasicAuthentication"
	separator = "/"
)

var (
	idleConnectionTimeout = 5 * time.Second
)

type Client struct {
	baseURL  *url.URL
	client   *http.Client
	auth     string
	username string
	password string
}

func NewClient(baseURL string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	u, _ := url.Parse(baseURL)
	return &Client{
		baseURL: u,
		client:  &http.Client{},
	}
}

func NewClientWithCerts(baseURL, certificatePath, key string, isInsecure bool) (*Client, error) {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	u, _ := url.Parse(baseURL)
	tlsConfig, err := GetTlsConfig(certificatePath, key, isInsecure)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL: u,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}, nil
}

func NewBasicAuthClient(baseURL string, user, pwd string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	u, _ := url.Parse(baseURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout: idleConnectionTimeout,
	}
	httpClient := &http.Client{
		Transport: tr,
	}
	return &Client{
		baseURL:  u,
		client:   httpClient,
		auth:     basicAuth,
		username: user,
		password: pwd,
	}
}

func NewBasicAuthClientWithProxy(baseURL string, user, pwd string) *Client {
	cli := NewBasicAuthClient(baseURL, user, pwd)
	cli.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout: idleConnectionTimeout,
		Proxy:           http.ProxyFromEnvironment,
	}
	return cli
}

func AppendPath(url, path string) string {
	if !strings.HasSuffix(url, separator) && !strings.HasPrefix(path, separator) {
		url = url + separator
	}
	return url + path
}

func (c *Client) NewGetRequest(ctx context.Context, apiPath string, body interface{}) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodGet, apiPath, body)
}

func (c *Client) NewPostRequest(ctx context.Context, apiPath string, body interface{}) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodPost, apiPath, body)
}

func (c *Client) newRequest(ctx context.Context, method, apiPath string, body interface{}) (*http.Request, error) {

	// Form the URL path
	u, _ := c.baseURL.Parse(apiPath)

	buf, err := encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set content type
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req = req.WithContext(ctx)
	return req, nil
}

func (c *Client) GetJson(ctx context.Context, apiPath string, jsonResponse interface{}) (*http.Response, error) {
	r, e := c.NewGetRequest(ctx, apiPath, nil)
	if e != nil {
		return nil, e
	}
	return c.DoJson(r, jsonResponse)
}

func (c *Client) DoJson(r *http.Request, jsonResponse interface{}) (*http.Response, error) {
	resp, err := c.do(r)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)
	if jsonResponse != nil {
		err = json.NewDecoder(resp.Body).Decode(jsonResponse)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (c *Client) Get(r *http.Request) (*http.Response, error) {
	resp, err := c.do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DoStr(r *http.Request) (string, error) {
	resp, err := c.do(r)
	if err != nil {
		return "", err
	}
	defer closeResponseBody(resp)
	buf := new(bytes.Buffer)
	if c.IsErrResponse(resp.StatusCode) {
		_, err = buf.ReadFrom(resp.Body)
		return "", err
	}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) IsErrResponse(code int) bool {
	return !(code >= 200 && code <= 399)
}

func (c *Client) GetStr(ctx context.Context, apiPath string) (string, error) {
	r, e := c.NewGetRequest(ctx, apiPath, nil)
	if e != nil {
		return "", e
	}
	return c.DoStr(r)
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	if c.auth == basicAuth {
		r.SetBasicAuth(c.username, c.password)
	}
	return c.client.Do(r)
}

func encode(v interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter
	if v != nil {
		buf = new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func setCookie() *cookiejar.Jar {
	var cookies []*http.Cookie
	jar, _ := cookiejar.New(nil)
	cookie := &http.Cookie{
		Name:  "RetryOnServerError",
		Value: "true",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse("http://localhost")
	jar.SetCookies(u, cookies)
	return jar
}

func isFileExists(fileName string) bool {
	isFileExists, _ := IsFileExists(fileName)
	return isFileExists
}

func IsUrlExists(web string) bool {
	response, errors := http.Get(web)
	if response != nil && response.Body != nil {
		defer closeResponseBody(response)
	}
	if errors != nil {
		return false
	}
	if response.StatusCode == 200 {
		return true
	}
	return false
}

func closeResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Printf("unable to close response body")
	}
}
