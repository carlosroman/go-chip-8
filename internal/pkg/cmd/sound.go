package cmd

import (
	"github.com/hajimehoshi/oto"
	log "github.com/sirupsen/logrus"
	"math"
)

const (
	sampleRate = 44100
	//bufferSize = 2048
	//bufferSize = 735
	bufferSize = 1470
	freq       = 587.3 // freqD
)

type soundCard struct {
	player *oto.Player
	sample []byte
}

func (s *soundCard) processSound(soundChan <-chan byte) (err error) {
	go func() {}()
	for b := range soundChan {
		if b > 0x0 {
			if _, err = s.player.Write(s.sample); err != nil {
				return err
			}
		}
	}
	return err
}

func newSoundCard() (s *soundCard, err error) {
	p, err := oto.NewPlayer(sampleRate, 1, 2, bufferSize)
	if err != nil {
		log.WithError(err).Warn("Failed creating sound card")
		return s, err
	}

	const length = sampleRate / 60
	out := make([]int16, length)
	const vol = 1.0 / 16.0
	square(out, vol, freq, 0.25)
	bytes := make([]byte, length*2)
	for i := range out {
		bytes[2*i] = byte(out[i])
		bytes[2*i+1] = byte(out[i] >> 8)
	}
	s = &soundCard{
		player: p,
		sample: bytes,
	}
	return s, err
}

func square(out []int16, volume float64, freq float64, sequence float64) {
	if freq == 0 {
		for i := 0; i < len(out); i++ {
			out[i] = 0
		}
		return
	}
	length := int(float64(sampleRate) / freq)
	if length == 0 {
		panic("invalid freq")
	}
	for i := 0; i < len(out); i++ {
		a := int16(volume * math.MaxInt16)
		if i%length < int(float64(length)*sequence) {
			a = -a
		}
		out[i] = a
	}
}
