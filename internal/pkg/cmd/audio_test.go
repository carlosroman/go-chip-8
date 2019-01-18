package cmd

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

func TestSoundCard(t *testing.T) {
	b := []byte{0x0, 0x1, 0x2}
	w := &mockWriter{}
	s := newSoundCard(b, w)
	w.On("Write", b).Return(1, nil)
	var wg sync.WaitGroup
	wg.Add(1)
	sc := make(chan byte, 2)
	go func() {
		defer wg.Done()
		err := s.ProcessSound(sc)
		assert.NoError(t, err)
	}()
	sc <- 0x0
	sc <- 0x1
	close(sc)
	wg.Wait()
	w.AssertNumberOfCalls(t, "Write", 1)
	w.AssertExpectations(t)
}

func TestSoundCard_error(t *testing.T) {
	b := []byte{0x0, 0x1, 0x2}
	w := &mockWriter{}
	s := newSoundCard(b, w)
	exp := errors.New("something went wrong")
	w.On("Write", b).Return(0, exp)
	var wg sync.WaitGroup
	wg.Add(1)
	sc := make(chan byte, 2)
	go func() {
		defer wg.Done()
		err := s.ProcessSound(sc)
		assert.Error(t, err, exp.Error())
	}()
	sc <- 0x0
	sc <- 0x1
	close(sc)
	wg.Wait()
	w.AssertNumberOfCalls(t, "Write", 1)
	w.AssertExpectations(t)
}

func TestGenerateSample(t *testing.T) {
	a := generateSample()
	assert.Len(t, a, bufferSize)
	// TODO: Find a nice way to validate `a`
}

type mockWriter struct {
	mock.Mock
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}
