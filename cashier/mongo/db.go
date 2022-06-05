package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionCashier = "cashiers"
)

type Store interface {
	Collection() *mongo.Collection
	Name() string
}

type store struct {
	store *mongo.Collection
	name  string
}

func NewStore(ctx context.Context, cfg Config) (Store, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseConnString()))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	if err = conn.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("unable to ping db: %v", err)
	}

	return store{
		store: conn.Database(cfg.DatabaseName()).Collection(CollectionCashier),
		name:  cfg.DatabaseName(),
	}, nil
}

func (s store) Collection() *mongo.Collection {
	return s.store
}

func (s store) Name() string {
	return s.name
}
