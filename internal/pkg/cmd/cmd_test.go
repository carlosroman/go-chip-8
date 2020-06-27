package cmd

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/carlosroman/go-chip-8/pkg/cpu"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	bcChip8TestPath = "../../../test/roms/BC_test.ch8"
)

func TestGetCommand(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := GetCommand(ctx, &noopScreen{}, cpu.NewKeyboard(), &ctxLoop{}, func() (ap AudioPlayer, err error) {
		m := mockAudioPlayer{}
		m.On("ProcessSound", mock.Anything).Return(nil)
		return &m, nil
	})
	c.SetArgs([]string{"--rom", bcChip8TestPath})
	var wg sync.WaitGroup
	wg.Add(2)
	var err error
	go func(w *sync.WaitGroup) {
		defer w.Done()
		log.Info("Starting app")
		_, err = c.ExecuteC()
		log.WithError(err).Info("Exited app")
		assert.NoError(t, err)
	}(&wg)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		<-time.After(300 * time.Millisecond) // allow about 18 ticks
		cancel()
	}(&wg)
	wg.Wait()
}

type noopScreen struct {
}

func (s *noopScreen) Draw(frameBuffer []byte) {

}

type ctxLoop struct {
	err error
}

func (l *ctxLoop) Run(ctx context.Context) error {
	<-ctx.Done()
	return l.err
}

type mockAudioPlayer struct {
	mock.Mock
}

func (m *mockAudioPlayer) ProcessSound(soundChan <-chan byte) (err error) {
	args := m.Called(soundChan)
	for i := range soundChan {
		log.Debug(i)
	}
	return args.Error(0)
}
