package cpu

import (
	"github.com/stretchr/testify/assert"
	"sync"

	log "github.com/sirupsen/logrus"
	"testing"
)

func TestKeyboard_isKeyPressed(t *testing.T) {
	k := newKeyboard()
	k.keyPressed(0xa)
	assert.True(t, k.isKeyPressed(0xa))
}

func TestKeyboard_waitForKeyPressed(t *testing.T) {
	k := newKeyboard()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Wait()
		assert.False(t, k.isKeyPressed(0xa))
		log.Info("pressing key")
		k.keyPressed(0xb)
	}()
	wg.Done()
	log.Info("waiting for key")
	key := k.waitForKeyPressed()
	log.Info("checking key")
	assert.True(t, k.isKeyPressed(0xb))
	assert.Equal(t, byte(0xb), key)
}
