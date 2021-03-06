package main

import (
	"context"
	"log"

	pbt "github.com/murdho/playlists-by-tallinn"
	"github.com/murdho/playlists-by-tallinn/internal"
	"github.com/murdho/playlists-by-tallinn/internal/logger"
)

const (
	// currentTrack = ""
	// persists     = false
	currentTrack = "La La - Land Yo"
	persists     = true
)

func main() {
	pbt.InitSystem(
		new(testRadio),
		new(testStorage),
		logger.New(logger.DebugLevel),
	)

	if err := pbt.PlaylistsByTallinn(context.Background(), pbt.PubSubMessage{}); err != nil {
		log.Fatal(err)
	}
}

type testRadio struct{}

func (tr *testRadio) CurrentTrack() (string, error) {
	return currentTrack, nil
}

type testStorage struct{}

func (ts *testStorage) LoadTrack(ctx context.Context, trackName string) (*internal.Track, error) {
	track := &internal.Track{
		Name:     currentTrack,
		Persists: persists,
	}

	return track, nil
}

func (ts *testStorage) SaveTrack(ctx context.Context, track *internal.Track) error {
	return nil
}
