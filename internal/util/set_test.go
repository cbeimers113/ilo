package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	var1 := "hello"
	var2 := "world"

	// Test constructor without items
	set := NewSet()
	assert.False(t, set.Contains(var1))

	// Test constructor with items
	set = NewSet(var1, var2)
	assert.True(t, set.Contains(var1))
	assert.True(t, set.Contains(var2))

	// Test insert
	set = NewSet()
	assert.False(t, set.Contains(var1))
	set.Insert(var1)
	assert.True(t, set.Contains(var1))
}
