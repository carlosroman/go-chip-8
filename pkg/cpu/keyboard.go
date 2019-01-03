package cpu

import "sync"

type Keyboard interface {
	waitForKeyPressed() (key byte)
	isKeyPressed(key byte) bool
	keyPressed(key byte)
}

type keyboard struct {
	loc  sync.RWMutex
	kp   byte
	wait chan byte
}

func newKeyboard() Keyboard {
	return &keyboard{
		wait: make(chan byte, 1),
	}
}

func (k *keyboard) isKeyPressed(key byte) bool {
	k.loc.RLock()
	defer k.loc.RUnlock()
	return k.kp == key
}

func (k *keyboard) waitForKeyPressed() (key byte) {
	key = <-k.wait
	return key
}

func (k *keyboard) keyPressed(key byte) {
	k.loc.Lock()
	defer k.loc.Unlock()
	k.kp = key
	k.wait <- key
}
