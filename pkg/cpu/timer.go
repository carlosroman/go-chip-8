package cpu

import "sync"

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
