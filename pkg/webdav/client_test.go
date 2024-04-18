package webdav

import (
	"context"
	"fmt"
	"testing"
)

// path
func TestWebdav(t *testing.T) {
	client, _ := CreateClient(CreatHttpClientWithAuth(&httpClient{}, "", ""), "")
	fmt.Println(client.ReadDir(context.Background(), "", false))
}
