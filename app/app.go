package app

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/deuterium34/go-rave/player"
	"github.com/deuterium34/go-rave/visual"

	"github.com/alexflint/go-arg"
	"github.com/eiannone/keyboard"
)

type App struct {
	Args   AppArgs
	Player *player.Player
	Visual *visual.Visual

	CloseCh   chan (error)
	closed    atomic.Bool
	trackFile *os.File
}

type AppArgs struct {
	Input      string `arg:"positional,required"`
	Mono       bool   `arg:"-m,--mono" default:"false"`
	SampleRate int    `arg:"-s,--sample-rate" defaul:"44100" help:"use 44100 (default) or 48000"`
}

func (AppArgs) Version() string {
	return "go-rave v0.1\n\ngithub.com/deuterium34"
}

var errClosed error = errors.New("App closed")

func NewApp() (*App, error) {
	var args AppArgs
	arg.MustParse(&args)
	if args.SampleRate == 0 {
		args.SampleRate = 44100
	}

	err := keyboard.Open()
	if err != nil {
		return nil, fmt.Errorf("keyboard.Open: %w", err)
	}

	pl, err := player.NewPlayer(args.Mono, args.SampleRate)
	if err != nil {
		return nil, fmt.Errorf("player.NewPlayer: %w", err)
	}

	vs := visual.NewVisual()

	app := &App{
		Player:  pl,
		Visual:  vs,
		Args:    args,
		CloseCh: make(chan error, 1),
	}

	return app, nil
}

func (a *App) Start() error {
	if a.closed.Load() {
		return errClosed
	}

	trackReader, tags, err := a.processFile(a.Args.Input)
	if err != nil {
		return fmt.Errorf("a.processFile: %w", err)
	}

	err = a.loadPlay(trackReader, tags)
	if err != nil {
		return fmt.Errorf("a.loadPlay: %w", err)
	}

	a.display(tags)

	// start keyboard listen
	go a.keyboardEventsLoop()
	return nil
}

func (a *App) Close() {
	a.internalClose(nil)
}
