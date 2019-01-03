package cpu

import "sync"

type keyboard struct {
	loc  sync.RWMutex
	kp   byte
	wait chan bool
}

func newKeyboard() *keyboard {
	return &keyboard{
		wait: make(chan bool, 1),
	}
}

func (k *keyboard) isKeyPressed(key byte) bool {
	k.loc.RLock()
	defer k.loc.RUnlock()
	return k.kp == key
}

func (k *keyboard) waitForKeyPressed() {
	<-k.wait
}

func (k *keyboard) keyPressed(key byte) {
	k.loc.Lock()
	defer k.loc.Unlock()
	k.kp = key
	k.wait <- true
}
