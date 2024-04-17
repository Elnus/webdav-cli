package webdav

import (
	"net/http"

	"github.com/emersion/go-webdav"
)

var (
	username string //认证账号
	password string //认证密码
)

func CreateClient(c webdav.HTTPClient, endpoint string) (*webdav.Client, error) {
	return webdav.NewClient(c, endpoint)
}

type httpClient struct {
	hClient http.Client
}

var _ webdav.HTTPClient = &httpClient{}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.hClient.Do(req)
}

func CreatHttpClientWithAuth(webHC webdav.HTTPClient) webdav.HTTPClient {
	return webdav.HTTPClientWithBasicAuth(webHC, username, password)
}
