package noop

import (
	"context"
	"sync"
	"time"

	"github.com/carlosroman/go-chip-8/pkg/cpu"
	log "github.com/sirupsen/logrus"
)

type Screen struct {
	lock sync.Mutex
	fb   []byte
}

func (s *Screen) Draw(frameBuffer []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Info("Draw")
	copy(s.fb, frameBuffer)
}

func (s *Screen) Refresh() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Info("Refresh")
	return nil
}

func (s *Screen) Run(ctx context.Context) error {
	go func() {
		cpu.Start(ctx, "loop", time.Second, s.Refresh)
	}()
	<-ctx.Done()
	return nil
}
