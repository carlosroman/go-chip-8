package cmd

import (
	"context"
	"github.com/carlosroman/go-chip-8/pkg/cpu"
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	ScreenWidth          = 64
	ScreenHeight         = 32
	defaultSixtyHz       = 60
	defaultSOneHundredHz = 100
)

func GetCommand(ctx context.Context, screen cpu.Screen, keyboard cpu.Keyboard, loop Loop, getSoundCard func() (ap AudioPlayer, err error)) *cobra.Command {
	var romPath string
	timer := time.Second / defaultSixtyHz          // 60hz
	cpuClock := time.Second / defaultSOneHundredHz // 100hz
	runCmd := &cobra.Command{
		Use:   "chip8",
		Short: "Chip8 is a Chip 8 emulator",
		Long:  "Chip8 is a Chip 8 emulator",
		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sc := make(chan byte, 60)
			ti := cpu.NewTimer(sc)
			wg := sync.WaitGroup{}
			s, err := getSoundCard()
			if err != nil {
				log.WithError(err).Fatal("Could not create sound card")
			}
			wg.Add(4)

			go func(w *sync.WaitGroup) {
				defer w.Done()
				if err != nil {
					log.WithError(err).Fatal("Could not create sound card")
				}
				go func() {
					<-ctx.Done()
					close(sc)
				}()
				if err = s.ProcessSound(sc); err != nil {
					log.WithError(err).Fatal("Sound card crashed")
				}
			}(&wg)

			go func(w *sync.WaitGroup) {
				defer w.Done()
				log.Warn("Starting loop")
				if err := loop.Run(ctx); err != nil {
					log.WithError(err).Fatal("Loop failed")
				}
				log.Warn("Stopping loop")
			}(&wg)
			go func(w *sync.WaitGroup) {
				defer wg.Done()
				log.Warn("Starting timer")
				ti.Start(ctx, timer)
				log.Warn("Stopping timer")
			}(&wg)
			go func(w *sync.WaitGroup) {
				defer w.Done()
				m := state.InitMemory()
				f, err := os.Open(romPath)
				if err != nil {
					log.WithError(err).Panicf("Could not open file '%s'", romPath)
				}
				err = m.LoadMemory(f)
				if err != nil {
					log.WithError(err).Panicf("Could not load memory with file '%s'", romPath)
				}
				s := rand.NewSource(time.Now().UnixNano())
				r := rand.New(s)
				c := cpu.NewCPU(m, r, keyboard, ti, screen)
				cpu.Start("cpu", ctx, cpuClock, c.Tick)
			}(&wg)
			wg.Wait()
		},
	}
	runCmd.Flags().StringVarP(&romPath, "rom", "r", "", "Path of rom to load (required)")
	if err := runCmd.MarkFlagRequired("rom"); err != nil {
		log.WithError(err).Fatal("Could not create command.")
	}
	return runCmd
}

type Loop interface {
	Run(ctx context.Context) error
}

type AudioPlayer interface {
	ProcessSound(soundChan <-chan byte) (err error)
}
