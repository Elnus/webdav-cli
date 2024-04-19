package webdav

import (
	"context"
	"fmt"
	"testing"

	"github.com/emersion/go-webdav"
)

// path
func TestWebdav(t *testing.T) {
	//client, _ := CreateClient(CreatHttpClientWithAuth(&httpClient{}, "", ""), "")
	client := InitClient(&httpClient{}, "", "", "")

	// readDir
	fmt.Println(client.ReadDir(context.Background(), "", false))

	// move
	//client.Move(context.Background(), "/via/settings.txt", ".settings.txt", &webdav.MoveOptions{NoOverwrite: true})

	// copy
	err := client.Copy(context.Background(), "", "", &webdav.CopyOptions{})
	fmt.Println(err)
}
