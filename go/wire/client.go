package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Config struct {
	Debug bool
}

func NewDefaultConfig(debug bool) *Config {
	return &Config{Debug: debug}
}

type ConnectionInfo struct {
	ProjectID string
}

func NewClient(ctx context.Context, info *ConnectionInfo) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, info.ProjectID)
	return client, err
}
