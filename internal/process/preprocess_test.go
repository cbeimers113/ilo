package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Preprocess(t *testing.T) {
	data := "saluton al ĉiuj!"
	want := "saluton al cxiuj!"
	got := Preprocess(data)
	assert.Equal(t, want, got)
}
