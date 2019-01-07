package cpu

import "sync"

type Keyboard interface {
	WaitForKeyPressed() (key byte)
	IsKeyPressed(key byte) bool
	KeyPressed(key byte)
}

type keyboard struct {
	loc  sync.RWMutex
	kp   byte
	wait chan byte
}

func NewKeyboard() Keyboard {
	return &keyboard{
		wait: make(chan byte, 1),
	}
}

func (k *keyboard) IsKeyPressed(key byte) bool {
	k.loc.RLock()
	defer k.loc.RUnlock()
	return k.kp == key
}

func (k *keyboard) WaitForKeyPressed() (key byte) {
	key = <-k.wait
	return key
}

func (k *keyboard) KeyPressed(key byte) {
	k.loc.Lock()
	defer k.loc.Unlock()
	k.kp = key
	k.wait <- key
}
