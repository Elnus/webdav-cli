package webdav

import (
	"context"
	"fmt"
	"testing"
)

var (
	webdavServer string //webdav服务器地址
	dir          string //操作目录
)

// path
func TestWebdav(t *testing.T) {
	client, _ := CreateClient(CreatHttpClientWithAuth(&httpClient{}), webdavServer)
	fmt.Println(client.ReadDir(context.Background(), dir, false))
}
