package webdav

import (
	"log"
	"net/http"

	"github.com/emersion/go-webdav"
)

type httpClient struct {
	hClient http.Client
}

var _ webdav.HTTPClient = &httpClient{}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.hClient.Do(req)
}

func CreateClient(c webdav.HTTPClient, endpoint string) (*webdav.Client, error) {
	return webdav.NewClient(c, endpoint)
}

func CreatHttpClientWithAuth(webHC webdav.HTTPClient, usr, pwd string) webdav.HTTPClient {
	return webdav.HTTPClientWithBasicAuth(webHC, usr, pwd)
}

func InitClient(webHC webdav.HTTPClient, endpoint string, usr, pwd string) *webdav.Client {
	c, err := CreateClient(CreatHttpClientWithAuth(webHC, usr, pwd), endpoint)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
