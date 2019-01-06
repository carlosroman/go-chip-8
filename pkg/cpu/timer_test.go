package cpu

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

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

func TestTimer_Start(t *testing.T) {
	ti := NewTimer()
	ti.SetSound(0x0d) // 13
	ti.SetDelay(0x33) // 51
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		ti.Start(ctx, time.Second/60)
		wg.Done()
	}()
	go func() {
		<-time.After(300 * time.Millisecond) // allow about 18 ticks
		cancel()
	}()
	wg.Wait()
	assert.Equal(t, byte(0x0), ti.GetSound(), "should stop at zero")
	assert.InDelta(t, byte(0x21), ti.GetDelay(), 1, "should be around 33 after 18 ticks")
}
