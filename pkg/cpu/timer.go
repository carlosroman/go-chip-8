package cpu

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type timer struct {
	lock      sync.RWMutex
	delay     byte
	sound     byte
	soundChan chan<- byte
}

func NewTimer(soundChan chan<- byte) *timer {
	return &timer{
		delay:     0,
		sound:     0,
		soundChan: soundChan,
	}
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

func (t *timer) tick() (err error) {
	log.Debug("tick")
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.delay > 0 {
		t.delay -= 1
	}
	if t.sound > 0 {
		t.sound -= 1
	}
	t.soundChan <- t.sound
	return err
}

func (t *timer) Start(ctx context.Context, duration time.Duration) {
	Start("timer", ctx, duration, t.tick)
}

func Start(name string, ctx context.Context, d time.Duration, tick func() error) {
	limit := rate.Every(d)
	log.WithField("name", name).WithField("d", d).WithField("limit", limit).Info("Starting timer")
	limiter := rate.NewLimiter(limit, 1)
	for {
		err := limiter.Wait(ctx)
		if err != nil {
			log.WithField("name", name).WithError(err).Warn("Got an error, exiting")
			break
		}
		if err = tick(); err != nil {
			log.WithField("name", name).WithError(err).Warn("Got an error running tick, exiting")
			break
		}
	}
}
