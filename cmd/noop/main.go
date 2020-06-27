package main

import (
	"context"

	"github.com/carlosroman/go-chip-8/internal/pkg/cmd"
	"github.com/carlosroman/go-chip-8/internal/pkg/noop"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	s := &noop.Screen{}
	ctx, _ := context.WithCancel(context.Background())
	c := cmd.GetCommand(ctx, s, cpu.NewKeyboard(), s, func() (ap cmd.AudioPlayer, err error) {
		return &audioPlayer{}, nil
	})
	if err := c.Execute(); err != nil {
		log.WithError(err).Fatal("app crashed")
	}
}

type audioPlayer struct {
}

func (n *audioPlayer) ProcessSound(soundChan <-chan byte) (err error) {
	for i := range soundChan {
		log.Debug(i)
	}
	return err
}
