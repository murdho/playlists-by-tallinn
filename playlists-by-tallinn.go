package playlistsbytallinn

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/murdho/playlists-by-tallinn/internal/logger"
	"github.com/murdho/playlists-by-tallinn/radio"
	"github.com/murdho/playlists-by-tallinn/storage"
)

func init() {
	gcpProject := os.Getenv("GCP_PROJECT")

	logLevel := logger.InfoLevel
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.DebugLevel
	}

	InitSystem(
		radio.NewRaadioTallinn(),
		storage.NewFirestoreStorage(gcpProject, "playlists-by-tallinn"),
		logger.New(logLevel),
	)
}

func PlaylistsByTallinn(ctx context.Context, _ PubSubMessage) error {
	sys.logger.Debug("starting")

	trackName, err := sys.radio.CurrentTrack()
	if err != nil {
		return errors.Wrap(err, "getting current track failed")
	}

	sys.logger.Info("current track", zap.String("name", trackName))

	if trackName == "" {
		sys.logger.Debug("current track empty, all done")
		return nil
	}

	sys.logger.Debug("loading track from storage")

	track, err := sys.trackStorage.LoadTrack(ctx, trackName)
	if err != nil {
		return errors.Wrap(err, "loading track from storage failed")
	}

	sys.logger.Debug(
		"track from storage",
		zap.String("name", track.Name),
		zap.Bool("persists", track.Persists),
	)

	if track.Persists {
		sys.logger.Debug("track already persists, all done")
		return nil
	}

	track.Persists = true
	sys.logger.Debug("saving track to storage")

	if err := sys.trackStorage.SaveTrack(ctx, track); err != nil {
		return errors.Wrap(err, "saving track to storage failed")
	}

	sys.logger.Debug("track saved to storage, all done")

	return nil
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}
