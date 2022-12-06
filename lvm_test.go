package go_lvm

import (
	"os"
	"reflect"
	"testing"

	"github.com/masahiro331/go-lvm/types"
	"github.com/stretchr/testify/require"
)

func Test_parseMetadata(t *testing.T) {
	type args struct {
		testFile string
	}
	tests := []struct {
		name string
		args args
		want types.Metadata
	}{
		{
			name: "happy path",
			args: args{
				testFile: "testdata/metadata.txt",
			},
			want: types.Metadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.testFile)
			require.NoError(t, err)
			defer f.Close()

			got, err := parseMetadata(f)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMetadata() got = %v, want %v", got, tt.want)
			}
		})
	}
}
