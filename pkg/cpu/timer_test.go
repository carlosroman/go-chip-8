package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimer_SetDelay(t *testing.T) {
	ti := NewTimer()
	assert.NotNil(t, ti)
	ti.SetDelay(0xaf)
	assert.Equal(t, byte(0xaf), ti.GetDelay())
}

func TestTimer_SetSound(t *testing.T) {
	ti := NewTimer()
	assert.NotNil(t, ti)
	ti.SetSound(0xaa)
	assert.Equal(t, byte(0xaa), ti.GetSound())
}
