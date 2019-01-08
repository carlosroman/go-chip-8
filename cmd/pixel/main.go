package main

import (
	"context"
	"github.com/carlosroman/go-chip-8/internal/pkg/cmd"
	"github.com/carlosroman/go-chip-8/internal/pkg/pixel"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	"github.com/faiface/pixel/pixelgl"
	log "github.com/sirupsen/logrus"
)

func main() {
	//  PixelGL needs ro use the main thread for all the windowing and graphics code
	pixelgl.Run(run)
}

func run() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "15:04:05.999999"
	log.SetLevel(log.WarnLevel)
	log.SetFormatter(customFormatter)
	k := cpu.NewKeyboard()
	ctx, cancel := context.WithCancel(context.Background())
	s, err := pixel.NewScreen(cmd.ScreenWidth, cmd.ScreenHeight, cancel, k)
	if err != nil {
		log.WithError(err).Fatal("could not create a screen")
	}
	c := cmd.GetCommand(ctx, s, k, s)
	if err := c.Execute(); err != nil {
		log.WithError(err).Fatal("app crashed")
	}
}
