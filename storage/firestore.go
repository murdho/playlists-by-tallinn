package storage

import (
	"context"
	"crypto/md5"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/murdho/playlists-by-tallinn/internal"
	"github.com/murdho/playlists-by-tallinn/internal/lazyfirestore"
)

func NewFirestoreStorage(gcpProject, collectionName string) *firestoreStorage {
	return &firestoreStorage{
		getFirestoreClient: lazyfirestore.NewClientFunc(gcpProject),
		collectionName:     collectionName,
	}
}

type firestoreStorage struct {
	getFirestoreClient func() (*firestore.Client, error)
	collectionName     string
}

func (f *firestoreStorage) LoadTrack(ctx context.Context, trackName string) (*internal.Track, error) {
	track := internal.NewTrack(trackName, false)

	client, err := f.getFirestoreClient()
	if err != nil {
		return nil, errors.Wrap(err, "getting Firestore client failed")
	}

	snapshot, err := f.getDocument(client, trackName).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return track, nil
		}

		return nil, errors.Wrap(err, "getting track data snapshot failed")
	}

	if err := snapshot.DataTo(&track); err != nil {
		return nil, errors.Wrap(err, "building track from snapshot failed")
	}

	track.Persists = true

	return track, nil
}

func (f *firestoreStorage) SaveTrack(ctx context.Context, track *internal.Track) error {
	client, err := f.getFirestoreClient()
	if err != nil {
		return errors.Wrap(err, "getting Firestore client failed")
	}

	if _, err := f.getDocument(client, track.Name).Set(ctx, track); err != nil {
		return errors.Wrap(err, "saving track failed")
	}

	return nil

}

func (f *firestoreStorage) getDocument(client *firestore.Client, trackName string) *firestore.DocumentRef {
	documentID := firestoreDocumentID(trackName)
	return client.Collection(f.collectionName).Doc(documentID)
}

func firestoreDocumentID(trackName string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(trackName)))
}
