package ebiten

import (
	"context"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
	"image"
	"sync"
)

func NewScreen(width int, height int, cancel context.CancelFunc, keyboard cpu.Keyboard) *Screen {
	os, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
	return &Screen{
		offscreen:      os,
		offscreenImage: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

type Screen struct {
	lock           sync.RWMutex
	offscreen      *ebiten.Image
	offscreenImage *image.RGBA
}

func (s *Screen) Draw(frameBuffer []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Info("Draw")
	for px := range frameBuffer {
		var col uint8
		if frameBuffer[px] > 0x0 {
			col = 0xff
		} else {
			col = 0x00
		}

		s.offscreenImage.Pix[4*px] = col
		s.offscreenImage.Pix[4*px+1] = col
		s.offscreenImage.Pix[4*px+2] = col
		s.offscreenImage.Pix[4*px+3] = 0xff
	}
	_ = s.offscreen.ReplacePixels(s.offscreenImage.Pix)
}

func (s *Screen) Run(ctx context.Context) (err error) {
	w, h := s.offscreen.Size()
	err = ebiten.Run(s.update, w, h, 10, "CHIP-8")
	return err
}

func (s *Screen) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	_ = screen.DrawImage(s.offscreen, nil)
	return nil
}

func (s *Screen) ProcessSound(soundChan <-chan byte) (err error) {
	for i := range soundChan {
		log.Debug(i)
	}
	return nil
}
