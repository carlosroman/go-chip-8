package noop

import (
	"context"
	"github.com/carlosroman/go-chip-8/internal/pkg/cmd"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func NewScreen() *Screen {
	return &Screen{
		fb: make([]byte, cmd.ScreenWidth*cmd.ScreenHeight),
	}
}

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
	var screen [][]byte
	for y := 0; y < 32; y++ {
		row := make([]byte, 64)
		for x := 0; x < 64; x++ {
			px := x + (y * 64)
			row[x] = s.fb[px]
		}
		screen = append(screen, row)
	}
	//for y:=31;y>=0;y--{
	//	log.Info(screen[y])
	//}
	return nil
}

func (n *Screen) Run(ctx context.Context) error {
	go func() {
		cpu.Start("loop", ctx, time.Second, n.Refresh)
	}()
	<-ctx.Done()
	return nil
}
