package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestCheckDirOrFileIsExist(t *testing.T) {
	_, err := os.Stat("/tmp")
	// fmt.Println(os.IsExist(err))
	fmt.Println(os.IsNotExist(err))
}
