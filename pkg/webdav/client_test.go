package webdav

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/emersion/go-webdav"
)

// path
func TestRead(t *testing.T) {
	//client, _ := CreateClient(CreatHttpClientWithAuth(&httpClient{}, "", ""), "")
	client := InitClient(&httpClient{}, "", "", "")
	fmt.Println(client.ReadDir(context.Background(), "", false))
}

func TestCreate(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	file, _ := client.Create(context.Background(), "")
	byte, _ := os.ReadFile("")
	file.Write(byte)
	defer file.Close()
}

func TestFindCurrentUserPrincipal(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	fmt.Println(client.FindCurrentUserPrincipal(context.Background()))
}

func TestRemove(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	err := client.RemoveAll(context.Background(), "")
	fmt.Println(err)
}

func TestStat(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	fileInfo, _ := client.Stat(context.Background(), "")
	fmt.Println(fileInfo.Path)
}

func TestOpen(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	file, _ := client.Open(context.Background(), "")
	data, _ := io.ReadAll(file)
	defer file.Close()
	//fmt.Println(string(data))
	osFile, _ := os.Create(".")
	osFile.Write(data)
}

func TestCopy(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	err := client.Copy(context.Background(), "", "", &webdav.CopyOptions{})
	fmt.Println(err)
}

func TestMove(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	client.Move(context.Background(), "", "", &webdav.MoveOptions{NoOverwrite: true})
}

func TestMkdir(t *testing.T) {
	client := InitClient(&httpClient{}, "", "", "")
	err := client.Mkdir(context.Background(), "")
	fmt.Println(err)
}
