package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cbeimers113/ilo/internal/locale"
)

var good_test_cfg = []byte("locale: en")
var bad_test_cfg = []byte("this: field: errors")

func Test_Config(t *testing.T) {
	_, err := New(bad_test_cfg)
	assert.Error(t, err)

	cfg, err := New(good_test_cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "en", cfg.Locale)
}

func Test_Message(t *testing.T) {
	cfg, err := TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		index int
		want  string
	}{
		{
			name:  "happy path - get indexed message",
			index: locale.ErrNoArguments,
			want:  "no arguments supplied, nothing to do",
		},
		{
			name:  "sad path - get index out of bounds",
			index: 1000,
			want:  "no message at index 1000 for locale \"en\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.Message(tt.index)
			assert.Equal(t, tt.want, got)
		})
	}
}
