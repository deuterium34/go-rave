package app

import (
	"fmt"
	"io"
	"os"

	"github.com/deuterium34/go-rave/parser"

	"github.com/eiannone/keyboard"
)

func (a *App) processFile(filepath string) (io.Reader, parser.Tags, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, parser.UnknownTags(), fmt.Errorf("os.Open: %w", err)
	}
	a.trackFile = file

	trackR, tags, err := parser.Parse(a.trackFile)
	if err != nil {
		return nil, parser.UnknownTags(), fmt.Errorf("parser.Parse: %w", err)
	}

	if seeker, ok := trackR.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	return trackR, tags, nil
}

func (a *App) loadPlay(r io.Reader, tags parser.Tags) error {
	a.Player.LoadTrack(r)
	err := a.Player.Play()
	if err != nil {
		return fmt.Errorf("a.Player.Play: %w", err)
	}

	a.Visual.LoadTrack(tags.Length)
	err = a.Visual.Play()
	if err != nil {
		return fmt.Errorf("a.Visual.Play: %w", err)
	}

	return nil
}

func (a *App) display(tg parser.Tags) {
	fmt.Printf("Now playing: %s\n", a.Args.Input)
	fmt.Println(Guitar)
	fmt.Printf("%s - %s\n%s\n\n", tg.Author, tg.Title, tg.Album)
}

func (a *App) keyboardEventsLoop() {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			a.internalClose(err)
			return
		}

		if key == keyboard.KeyEsc {
			a.internalClose(nil)
			return
		}

		if key == keyboard.KeySpace || char == ' ' {
			a.Player.PlayPause()
			a.Visual.PlayPause()
		}
	}
}

func (a *App) internalClose(reason error) {
	if a.closed.Swap(true) {
		return
	}

	keyboard.Close()
	a.Player.Close()
	a.Visual.Close()
	if a.trackFile != nil {
		a.trackFile.Close()
	}

	a.CloseCh <- reason
}
