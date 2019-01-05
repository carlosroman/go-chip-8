package cpu

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type timer struct {
	lock  sync.RWMutex
	delay byte
	sound byte
}

func NewTimer() *timer {
	return &timer{}
}

func (t *timer) SetDelay(val byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.delay = val
}

func (t *timer) SetSound(val byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.sound = val
}

func (t *timer) GetDelay() (val byte) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.delay
}

func (t *timer) GetSound() (val byte) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.sound
}

func (t *timer) tick() {
	log.Debug("tick")
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.delay > 0 {
		t.delay -= 1
	}
	if t.sound > 0 {
		t.sound -= 1
	}
}

func (t *timer) Start(ctx context.Context) {
	limit := rate.Every(time.Second / 60)
	limiter := rate.NewLimiter(limit, 1)
	for {
		err := limiter.Wait(ctx)
		if err != nil {
			log.WithError(err).Info("Got an error, exiting")
			break
		}
		t.tick()
	}
}
