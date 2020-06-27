package cmd

import (
	"io"
	"math"

	"github.com/hajimehoshi/oto"
	log "github.com/sirupsen/logrus"
)

const (
	sampleRate = 44100
	bufferSize = 1470
	freq       = 587.3 // freqD
	sequence   = 0.25
	volume     = 0.0625 // 1.0 / 16.0
)

type soundCard struct {
	player io.Writer
	sample []byte
}

func (s *soundCard) ProcessSound(soundChan <-chan byte) (err error) {
	for b := range soundChan {
		if b > 0x0 {
			if _, err = s.player.Write(s.sample); err != nil {
				return err
			}
		}
	}
	return err
}

func newSoundCard(sample []byte, player io.Writer) (s *soundCard) {
	return &soundCard{
		player: player,
		sample: sample,
	}
}

func setupSoundCard() (s *soundCard, err error) {
	bytes := generateSample()
	p, err := oto.NewPlayer(sampleRate, 1, 2, bufferSize)
	if err != nil {
		log.WithError(err).Warn("Failed creating sound card")
		return s, err
	}
	return newSoundCard(bytes, p), err
}

func generateSample() []byte {
	const length = bufferSize / 2
	s := make([]int16, length)
	fill(s)
	bytes := make([]byte, bufferSize)
	for i := range s {
		bytes[2*i] = byte(s[i])
		bytes[2*i+1] = byte(s[i] >> 8)
	}
	return bytes
}

func fill(sample []int16) {
	vol := float64(volume)
	f := float64(sampleRate)
	length := int(f / float64(freq))
	for i := 0; i < len(sample); i++ {
		a := int16(vol * math.MaxInt16)
		if i%length < int(float64(length)*sequence) {
			a = -a
		}
		sample[i] = a
	}
}
