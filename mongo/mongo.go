package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds all client options for connecting to mongo
type Config struct {
	Timeout time.Duration `envconfig:"TIMEOUT"`
	Name    string        `envconfig:"NAME"`
	URI     string        `envconfig:"URI"`
}

// New instantiate mongo client based on configuration
func New(ctx context.Context, cfg Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	return client, nil
}
