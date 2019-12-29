package main

import (
	"context"
	"github.com/carlosroman/go-chip-8/internal/pkg/cmd"
	ebiten8 "github.com/carlosroman/go-chip-8/internal/pkg/ebiten"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	log "github.com/sirupsen/logrus"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "15:04:05.999999"
	log.SetLevel(log.WarnLevel)
	log.SetFormatter(customFormatter)
	k := cpu.NewKeyboard()
	ctx, cancel := context.WithCancel(context.Background())
	s := ebiten8.NewScreen(cmd.ScreenWidth, cmd.ScreenHeight, cancel, k)

	c := cmd.GetCommand(ctx, s, k, s, func() (ap cmd.AudioPlayer, err error) {
		return s, nil
	})
	if err := c.Execute(); err != nil {
		log.WithError(err).Fatal("app crashed")
	}
}
