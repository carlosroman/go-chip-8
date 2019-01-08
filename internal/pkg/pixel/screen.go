package pixel

import (
	"context"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"sync"
	"time"
)

func NewScreen(width int, height int, cancel context.CancelFunc, keyboard cpu.Keyboard) (s *Screen, err error) {
	cfg := pixelgl.WindowConfig{
		Title:  "CHIP-8",
		Bounds: pixel.R(0, 0, float64(width*10), float64(height*10)),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.WithError(err).Error("Could not create window")
		return nil, err
	}
	win.Clear(colornames.Black)
	return &Screen{
		fb:       make([]byte, 64*32),
		w:        float64(width),
		h:        float64(height),
		gridDraw: imdraw.New(nil),
		win:      win,
		cancel:   cancel,
		keyboard: keyboard,
	}, err
}

type Screen struct {
	lock     sync.Mutex
	fb       []byte
	w        float64
	h        float64
	win      *pixelgl.Window
	gridDraw *imdraw.IMDraw
	cancel   context.CancelFunc
	keyboard cpu.Keyboard
}

func (s *Screen) Draw(frameBuffer []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Info("Draw")
	copy(s.fb, frameBuffer)
}

func (s *Screen) refresh() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.win.Closed() {
		return errors.New("windows was closed")
	}

	if s.win.JustPressed(pixelgl.KeyEscape) {
		defer s.cancel()
		return nil
	}

	/**
	Keypad                   Keyboard
	+-+-+-+-+                +-+-+-+-+
	|1|2|3|C|                |1|2|3|4|
	+-+-+-+-+                +-+-+-+-+
	|4|5|6|D|                |Q|W|E|R|
	+-+-+-+-+       =>       +-+-+-+-+
	|7|8|9|E|                |A|S|D|F|
	+-+-+-+-+                +-+-+-+-+
	|A|0|B|F|                |Z|X|C|V|
	+-+-+-+-+                +-+-+-+-+
	*/

	if s.win.Pressed(pixelgl.Key1) { // ROW 0
		s.keyboard.KeyPressed(0x1)
	} else if s.win.Pressed(pixelgl.Key2) {
		s.keyboard.KeyPressed(0x2)
	} else if s.win.Pressed(pixelgl.Key3) {
		s.keyboard.KeyPressed(0x3)
	} else if s.win.Pressed(pixelgl.Key4) {
		s.keyboard.KeyPressed(0xC)
	} else if s.win.Pressed(pixelgl.KeyQ) { // ROW 1
		s.keyboard.KeyPressed(0x4)
	} else if s.win.Pressed(pixelgl.KeyW) {
		s.keyboard.KeyPressed(0x5)
	} else if s.win.Pressed(pixelgl.KeyE) {
		s.keyboard.KeyPressed(0x6)
	} else if s.win.Pressed(pixelgl.KeyR) {
		s.keyboard.KeyPressed(0xD)
	} else if s.win.Pressed(pixelgl.KeyA) { // ROW 1
		s.keyboard.KeyPressed(0x7)
	} else if s.win.Pressed(pixelgl.KeyS) {
		s.keyboard.KeyPressed(0x8)
	} else if s.win.Pressed(pixelgl.KeyD) {
		s.keyboard.KeyPressed(0x9)
	} else if s.win.Pressed(pixelgl.KeyF) {
		s.keyboard.KeyPressed(0xE)
	} else if s.win.Pressed(pixelgl.KeyZ) { // ROW 2
		s.keyboard.KeyPressed(0xA)
	} else if s.win.Pressed(pixelgl.KeyX) {
		s.keyboard.KeyPressed(0x0)
	} else if s.win.Pressed(pixelgl.KeyC) {
		s.keyboard.KeyPressed(0xB)
	} else if s.win.Pressed(pixelgl.KeyV) {
		s.keyboard.KeyPressed(0xF)
	} else {
		s.keyboard.Clear()
	}
	s.gridDraw.Clear()
	// do grid stuff
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			px := x + (y * 64)
			if s.fb[px] > 0x0 {
				s.gridDraw.Color = colornames.White
			} else {
				s.gridDraw.Color = colornames.Black
			}
			s.gridDraw.Push(
				pixel.V(float64(x*10), float64((31-y)*10)),     // 10, 10
				pixel.V(float64(x*10+9), float64((31-y)*10+9)), // 19, 19
			)
			s.gridDraw.Rectangle(0)
		}
	}

	s.gridDraw.Draw(s.win)
	s.win.Update()
	return nil
}

func (s *Screen) Run(ctx context.Context) (err error) {
	cpu.Start("screen", ctx, 33*time.Millisecond, s.refresh)
	return nil
}
