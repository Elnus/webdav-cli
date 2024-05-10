package cmd

import (
	"net/http"
	"testing"
	wb "webdav-cli/pkg/webdav"
)

func TestCmd(t *testing.T) {
	vars.Client = wb.InitClient(&http.Client{}, "", "", "")
	vars.recursive = true
}
