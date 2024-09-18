package process

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"cbeimers113/ilo/internal/config"
)

func Test_ParseArgs(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantSource string
		wantTarget string
	}{
		{
			name:    "no arguments",
			wantErr: true,
		},
		{
			name:       "source file but no target",
			args:       []string{"test/test.ilo"},
			wantSource: "test/test.ilo",
			wantTarget: "test/test",
		},
		{
			name:       "source file and single word target",
			args:       []string{"test/test.ilo", "target"},
			wantSource: "test/test.ilo",
			wantTarget: "target",
		},
		{
			name:       "source file and multi word target",
			args:       []string{"test/test.ilo", "hello", "world"},
			wantSource: "test/test.ilo",
			wantTarget: "helloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, target, err := ParseArgs(cfg, tt.args)
			assert.Equal(t, tt.wantSource, source)
			assert.Equal(t, tt.wantTarget, target)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_checkArgs(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name string
		args []string
		want error
	}{
		{
			name: "sad path - no args",
			args: []string{},
			want: errors.New("missing argument(s)"),
		},
		{
			name: "sad path - no ilo source file",
			args: []string{"source.go"},
			want: errors.New("source file not supplied"),
		},
		{
			name: "sad path - source file not first",
			args: []string{"target", "source.ilo"},
			want: errors.New("source file not supplied"),
		},
		{
			name: "sad path - invalid characters in args",
			args: []string{"source@file,"},
			want: errors.New("invalid characters in text \"source@file,\""),
		},
		{
			name: "sad path - source file doesn't exist",
			args: []string{"source.ilo"},
			want: errors.New("source file doesn't exist: \"source.ilo\""),
		},
		{
			name: "happy path - one source file",
			args: []string{"test/test.ilo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkArgs(cfg, tt.args)
			assert.Equal(t, tt.want, err)
		})
	}
}
