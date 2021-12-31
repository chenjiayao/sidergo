package validate

import (
	"testing"

	"github.com/chenjiayao/goredistraning/redis/rediserr"
)

func TestValidateSet(t *testing.T) {
	type args struct {
		args [][]byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "key value",
			args: args{
				args: [][]byte{[]byte("key"), []byte("value")},
			},
			want: nil,
		},

		{
			name: "key value ex not",
			args: args{
				args: [][]byte{[]byte("key"), []byte("value"), []byte("ex"), []byte("not")},
			},
			want: rediserr.NOT_INTEGER_ERROR,
		},

		{
			name: "key value px not",
			args: args{
				args: [][]byte{[]byte("key"), []byte("value"), []byte("px"), []byte("not")},
			},
			want: rediserr.NOT_INTEGER_ERROR,
		},

		{
			name: "key value ex 100 nx",
			args: args{
				args: [][]byte{[]byte("key"), []byte("value"), []byte("ex"), []byte("100"), []byte("nx")},
			},
			want: nil,
		},

		{
			name: "key value ex 100 px 100",
			args: args{
				args: [][]byte{[]byte("key"), []byte("value"), []byte("ex"), []byte("100"), []byte("px"), []byte("100")},
			},
			want: rediserr.SYNTAX_ERROR,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSet(tt.args.args)

			//got 和 want 都是nil
			if got == tt.want && got == nil && tt.want == nil {
				return
			}

			//got 和 want 都是 error，并且两个 error 一样
			if got != nil && tt.want != nil && got.Error() == tt.want.Error() {
				return
			}

			t.Errorf("ValidateSet() = %v, want %v", got, tt.want)
		})
	}
}

func TestValidateMSet(t *testing.T) {
	args := [][]byte{
		[]byte("key1"),
		[]byte("value1"),
		[]byte("key2"),
		[]byte("value2"),
	}
	err := ValidateMSet(args)
	if err != nil {
		t.Errorf("ValidateMSet(args) = %v, want nil", err)
	}
}
