package chksum

import (
	"encoding/hex"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_chkMD5(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "file1",
			args: args{data: func() []byte {
				f, _ := os.Open("")
				defer f.Close()
				data, _ := io.ReadAll(f)
				return data
			}()},
			want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chkMD5(tt.args.data)
			if strGot := hex.EncodeToString(got); !reflect.DeepEqual(strGot, tt.want) {
				t.Errorf("chkMD5() = %v, want %v", strGot, tt.want)
			}
		})
	}
}
