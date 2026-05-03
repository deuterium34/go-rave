package visual

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Track progress bar width
const Width = 50

var (
	ErrTrackNotStarted = errors.New("track length not set")
	ErrClosed          = errors.New("visuals closed")
)

type Visual struct {
	ProgressTrackWidth int
	TrackTotalLen      time.Duration

	mu      sync.RWMutex
	current time.Duration
	paused  bool
	cancel  context.CancelFunc
	running bool
}

func NewVisual() *Visual {
	return &Visual{
		ProgressTrackWidth: Width,
		TrackTotalLen:      0,
		paused:             true,
	}
}

func (v *Visual) LoadTrack(trackLen time.Duration) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.running {
		v.stopLoop()
	}

	v.TrackTotalLen = trackLen
	v.current = 0
	v.paused = true

	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel
	v.running = true

	go v.renderLoop(ctx)

	return nil
}

func (v *Visual) Play() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.TrackTotalLen <= 0 {
		return ErrTrackNotStarted
	}
	v.paused = false
	return nil
}

func (v *Visual) Pause() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.paused = true
	return nil
}

func (v *Visual) PlayPause() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.TrackTotalLen <= 0 {
		return ErrTrackNotStarted
	}
	v.paused = !v.paused
	return nil
}

func (v *Visual) Close() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.stopLoop()
}

func (v *Visual) stopLoop() {
	if v.cancel != nil {
		v.cancel()
		v.running = false
	}
}

func (v *Visual) renderLoop(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v.mu.Lock()
			if !v.paused {
				v.current += 100 * time.Millisecond
				v.renderProgressTrack()

				if v.current >= v.TrackTotalLen {
					v.running = false
					v.mu.Unlock()
					fmt.Println("\nFinished")
					return
				}
			}
			v.mu.Unlock()
		}
	}
}

func (v *Visual) renderProgressTrack() {
	ratio := float64(v.current) / float64(v.TrackTotalLen)
	if ratio > 1 {
		ratio = 1
	}

	filledLength := int(float64(v.ProgressTrackWidth) * ratio)

	bar := strings.Repeat("=", filledLength) + strings.Repeat("-", v.ProgressTrackWidth-filledLength)

	currMM, currSS := int(v.current.Minutes()), int(v.current.Seconds())%60
	totMM, totSS := int(v.TrackTotalLen.Minutes()), int(v.TrackTotalLen.Seconds())%60

	fmt.Printf("\r%02d:%02d [%s] %02d:%02d\033[K", currMM, currSS, bar, totMM, totSS)
}
