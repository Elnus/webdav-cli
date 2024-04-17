package webdav

import (
	"net/http"

	"github.com/emersion/go-webdav"
)

func CreateClient(c webdav.HTTPClient, endpoint string) (*webdav.Client, error) {
	return webdav.NewClient(c, endpoint)
}

type httpClient struct {
}

var hc webdav.HTTPClient = &httpClient{}

func (*httpClient) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func Creat(hc webdav.HTTPClient) webdav.HTTPClient {
	hc.Do(nil)
	return webdav.HTTPClientWithBasicAuth(hc, "username", "password")
}
