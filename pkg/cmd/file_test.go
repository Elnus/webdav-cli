package cmd

import (
	"context"
	"net/http"
	"testing"
	wb "webdav-cli/pkg/webdav"
)

func TestCmd(t *testing.T) {
	vars.Client = wb.InitClient(&http.Client{}, "", "", "")
	vars.recursive = true
	downloadFunc(context.Background(), "", "")
}
