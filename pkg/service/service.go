package service

import (
	"context"
	"path"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"

	"github.com/go-kit/kit/log"
)

const (
	CollectionProducts = "products"
	CollectionDraws    = "draws"
)

type Service interface {
	ListDraws(ctx context.Context, productID string) ([]*Draw, error)
	GetDraw(ctx context.Context, productID string, drawID string) (*Draw, error)
}

func New(client *firestore.Client, logger log.Logger) Service {
	var svc Service
	svc = &service{client}
	svc = LoggingMiddleware(logger)(svc)
	return svc
}

type service struct {
	firestoreClient *firestore.Client
}

func (s *service) ListDraws(ctx context.Context, productID string) ([]*Draw, error) {
	p := path.Join(CollectionProducts, productID, CollectionDraws)
	iter := s.firestoreClient.Collection(p).Documents(ctx)
	defer iter.Stop()
	draws := make([]*Draw, 0)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var d Draw
		if err := dsnap.DataTo(&d); err != nil {
			return nil, err
		}
		draws = append(draws, &d)
	}
	return draws, nil
}

func (s *service) GetDraw(ctx context.Context, productID string, drawID string) (*Draw, error) {
	p := path.Join(CollectionProducts, productID, CollectionDraws, drawID)
	dsnap, err := s.firestoreClient.Doc(p).Get(ctx)
	if err != nil {
		return nil, err
	}
	var d Draw
	if err := dsnap.DataTo(&d); err != nil {
		return nil, err
	}
	return &d, nil
}
