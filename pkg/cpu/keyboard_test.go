package cpu

import (
	"github.com/stretchr/testify/assert"
	"sync"

	log "github.com/sirupsen/logrus"
	"testing"
)

func TestKeyboard_isKeyPressed(t *testing.T) {
	k := NewKeyboard()
	k.KeyPressed(0xa)
	assert.True(t, k.IsKeyPressed(0xa))
}

func TestKeyboard_waitForKeyPressed(t *testing.T) {
	k := NewKeyboard()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Wait()
		assert.False(t, k.IsKeyPressed(0xa))
		log.Info("pressing key")
		k.KeyPressed(0xb)
	}()
	wg.Done()
	log.Info("waiting for key")
	key := k.WaitForKeyPressed()
	log.Info("checking key")
	assert.True(t, k.IsKeyPressed(0xb))
	assert.Equal(t, byte(0xb), key)
}
