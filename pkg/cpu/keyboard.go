package cpu

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Keyboard interface {
	WaitForKeyPressed() (key byte)
	IsKeyPressed(key byte) bool
	KeyPressed(key byte)
	Clear()
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
	if log.IsLevelEnabled(log.DebugLevel) {
		log.
			WithField("key", key).
			WithField("kp", k.kp).
			Debug("IsKeyPressed")
	}
	return k.kp == key
}

func (k *keyboard) WaitForKeyPressed() (key byte) {
	log.Info("Waiting...")
	key = <-k.wait
	if log.IsLevelEnabled(log.DebugLevel) {
		log.
			WithField("key", key).
			WithField("kp", k.kp).
			Debug("WaitForKeyPressed")
	}
	return key
}

func (k *keyboard) KeyPressed(key byte) {
	k.loc.Lock()
	defer k.loc.Unlock()
	k.kp = key
	if len(k.wait) > 0 {
		<-k.wait
	}
	k.wait <- key
}

func (k *keyboard) Clear() {
	k.loc.Lock()
	defer k.loc.Unlock()
	if len(k.wait) > 0 {
		<-k.wait
	}
	k.kp = 0x11
}
