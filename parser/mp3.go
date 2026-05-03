package parser

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
)

type MP3parser struct {
}

func NewMP3parser() TrackParser {
	return &MP3parser{}
}

func (p *MP3parser) ParseTrack(file *os.File) (io.Reader, Tags, error) {
	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, Tags{}, fmt.Errorf("mp3.NewDecoder: %w", err)
	}

	var tags Tags

	m, err := tag.ReadFrom(file)
	if err == nil {
		tags = Tags{
			Title:  m.Title(),
			Author: m.Artist(),
			Album:  m.Album(),
		}
	} else {
		fmt.Printf(".mp3 tags parse error: %v\n", err)
		tags = UnknownTags()
	}

	totalBytes := decodedMp3.Length()
	sampleRate := decodedMp3.SampleRate()
	seconds := float64(totalBytes) / 4 / float64(sampleRate)
	totalTime := time.Duration(seconds * float64(time.Second))

	tags.Length = totalTime

	return decodedMp3, tags, nil
}
