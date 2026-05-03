package parser

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Tags struct {
	Title  string
	Author string
	Album  string
	Length time.Duration
}

type TrackParser interface {
	ParseTrack(file *os.File) (io.Reader, Tags, error)
}

var ErrUnknownFormat error = errors.New("Unknown track format")

func UnknownTags() Tags {
	return Tags{
		Title:  "Unknown",
		Author: "Unknown",
		Album:  "Unknown",
		Length: 0,
	}
}

func Parse(file *os.File) (io.Reader, Tags, error) {
	fileExt := filepath.Ext(file.Name())

	switch fileExt {
	case ".mp3":
		return NewMP3parser().ParseTrack(file)
	default:
		return nil, UnknownTags(), ErrUnknownFormat
	}
}
