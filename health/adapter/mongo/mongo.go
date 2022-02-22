package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo ...
type Mongo struct {
	name   string
	ctx    context.Context
	client *mongo.Client
}

// New ...
func New(ctx context.Context, client *mongo.Client) *Mongo {
	return &Mongo{ctx: ctx, client: client}
}

// SetName ...
func (m *Mongo) SetName(name string) {
	m.name = name
}

// Name ...
func (m *Mongo) Name() string {
	if len(m.name) == 0 {
		return "MONGO"
	}
	return m.name
}

// Check ...
func (m *Mongo) Check() error {
	return m.client.Ping(m.ctx, readpref.Primary())
}
