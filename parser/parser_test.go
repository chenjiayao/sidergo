package parser

import (
	"bytes"
	"testing"

	"github.com/chenjiayao/goredistraning/redis"
)

func TestParseFromSocket(t *testing.T) {

	var buf bytes.Buffer
	buf.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"))
	ch := make(chan redis.RedisRequet)
	go ParseFromSocket(&buf, ch)

	r := <-ch
	if r.ToStrings() != "SET key value" {
		t.Errorf("err: %s", r.ToStrings())
	}
}

func Test_parseCmdArgsCount(t *testing.T) {
	type args struct {
		header []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "*3\r\n",
			args: args{
				header: []byte("*3\r\n"),
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "*a\r\n",
			args: args{
				header: []byte("*a\r\n"),
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCmdArgsCount(tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCmdArgsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseCmdArgsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOneCmdArgsLen(t *testing.T) {
	type args struct {
		cmd []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "$3\r\n",
			args: args{
				cmd: []byte("$3\r\n"),
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "$r\r\n",
			args: args{
				cmd: []byte("$r\r\n"),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOneCmdArgsLen(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOneCmdArgsLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseOneCmdArgsLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
